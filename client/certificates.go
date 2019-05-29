package client

import (
	"errors"
	"fmt"
	"time"

	tls "github.com/nirmata/kyverno/pkg/tls"
	certificates "k8s.io/api/certificates/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

// Issues TLS certificate for webhook server using given PEM private key
// Returns signed and approved TLS certificate in PEM format
func (c *Client) GenerateTlsPemPair(props tls.TlsCertificateProps) (*tls.TlsPemPair, error) {
	privateKey, err := tls.TlsGeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	certRequest, err := tls.TlsCertificateGenerateRequest(privateKey, props)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to create certificate request: %v", err))
	}

	certRequest, err = c.submitAndApproveCertificateRequest(certRequest)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to submit and approve certificate request: %v", err))
	}

	tlsCert, err := c.fetchCertificateFromRequest(certRequest, 10)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to configure a certificate for the Kyverno controller. A CA certificate is required to allow the Kubernetes API Server to communicate with Kyverno. You can either provide a certificate or configure your cluster to allow certificate signing. Please refer to https://github.com/nirmata/kyverno/installation.md.: %v", err))
	}

	return &tls.TlsPemPair{
		Certificate: tlsCert,
		PrivateKey:  tls.TlsPrivateKeyToPem(privateKey),
	}, nil
}

// Submits and approves certificate request, returns request which need to be fetched
func (c *Client) submitAndApproveCertificateRequest(req *certificates.CertificateSigningRequest) (*certificates.CertificateSigningRequest, error) {
	certClient, err := c.GetCSRInterface()
	if err != nil {
		return nil, err
	}
	csrList, err := c.ListResource(CSRs, "")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to list existing certificate requests: %v", err))
	}

	for _, csr := range csrList.Items {
		if csr.GetName() == req.ObjectMeta.Name {
			err := c.DeleteResouce(CSRs, "", csr.GetName())
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Unable to delete existing certificate request: %v", err))
			}
			c.logger.Printf("Old certificate request is deleted")
			break
		}
	}

	unstrRes, err := c.CreateResource(CSRs, "", req)
	if err != nil {
		return nil, err
	}
	c.logger.Printf("Certificate request %s is created", unstrRes.GetName())

	res, err := convertToCSR(unstrRes)
	if err != nil {
		return nil, err
	}
	res.Status.Conditions = append(res.Status.Conditions, certificates.CertificateSigningRequestCondition{
		Type:    certificates.CertificateApproved,
		Reason:  "NKP-Approve",
		Message: "This CSR was approved by Nirmata kyverno controller",
	})
	res, err = certClient.UpdateApproval(res)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to approve certificate request: %v", err))
	}
	c.logger.Printf("Certificate request %s is approved", res.ObjectMeta.Name)

	return res, nil
}

// Fetches certificate from given request. Tries to obtain certificate for maxWaitSeconds
func (c *Client) fetchCertificateFromRequest(req *certificates.CertificateSigningRequest, maxWaitSeconds uint8) ([]byte, error) {
	// TODO: react of SIGINT and SIGTERM
	timeStart := time.Now()
	for time.Now().Sub(timeStart) < time.Duration(maxWaitSeconds)*time.Second {
		unstrR, err := c.GetResource(CSRs, "", req.ObjectMeta.Name)
		if err != nil {
			return nil, err
		}
		r, err := convertToCSR(unstrR)
		if err != nil {
			return nil, err
		}

		if r.Status.Certificate != nil {
			return r.Status.Certificate, nil
		}

		for _, condition := range r.Status.Conditions {
			if condition.Type == certificates.CertificateDenied {
				return nil, errors.New(condition.String())
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Cerificate fetch timeout is reached: %d seconds", maxWaitSeconds))
}

const (
	ns             string = "kyverno"
	tlskeypair     string = "tls.kyverno"
	tlskeypaircert string = "tls.crt"
	tlskeypairkey  string = "tls.key"
	tlsca          string = "tls-ca"
	tlscarootca    string = "rootCA.crt"
)

// CheckPrePreqSelfSignedCert checks if the required secrets are defined
// if the user is providing self-signed certificates,key pair and CA
func (c *Client) CheckPrePreqSelfSignedCert() bool {
	// Check if secrets are defined if user is specifiying self-signed certificates
	tlspairfound := true
	tlscafound := true
	_, err := c.GetResource(Secrets, ns, tlskeypair)
	if err != nil {
		tlspairfound = false
	}
	_, err = c.GetResource(Secrets, ns, tlsca)
	if err != nil {
		tlscafound = false
	}
	if tlspairfound == tlscafound {
		return true
	}
	// Fail if only one of them is defined
	c.logger.Printf("while using self-signed certificates specify both secrets %s/%s & %s/%s for (cert,key) pair & CA respectively", ns, tlskeypair, ns, tlsca)

	if !tlspairfound {
		c.logger.Printf("secret %s/%s not defined for (cert,key) pair", ns, tlskeypair)
	}

	if !tlscafound {
		c.logger.Printf("secret %s/%s not defined for CA", ns, tlsca)
	}
	return false
}

func (c *Client) TlsrootCAfromSecret() (result []byte) {
	stlsca, err := c.GetResource(Secrets, ns, tlsca)
	if err != nil {
		return result
	}
	tlsca, err := convertToSecret(stlsca)
	if err != nil {
		utilruntime.HandleError(err)
		return result
	}

	result = tlsca.Data[tlscarootca]
	if len(result) == 0 {
		c.logger.Printf("root CA certificate not found in secret %s/%s", ns, tlsca.Name)
		return result
	}
	c.logger.Printf("using CA bundle defined in secret %s/%s to validate the webhook's server certificate", ns, tlsca.Name)
	return result
}

func (c *Client) TlsPairFromSecrets() *tls.TlsPemPair {
	// Check if secrets are defined
	stlskeypair, err := c.GetResource(Secrets, ns, tlskeypair)
	if err != nil {
		return nil
	}
	tlskeypair, err := convertToSecret(stlskeypair)
	if err != nil {
		utilruntime.HandleError(err)
		return nil
	}

	pemPair := tls.TlsPemPair{
		Certificate: tlskeypair.Data[tlskeypaircert],
		PrivateKey:  tlskeypair.Data[tlskeypairkey],
	}

	if len(pemPair.Certificate) == 0 {
		c.logger.Printf("TLS Certificate not found in secret %s/%s", ns, tlskeypair.Name)
		return nil
	}
	if len(pemPair.PrivateKey) == 0 {
		c.logger.Printf("TLS PrivateKey not found in secret %s/%s", ns, tlskeypair.Name)
		return nil
	}
	c.logger.Printf("using TLS pair defined in secret %s/%s for webhook's server tls configuration", ns, tlskeypair.Name)
	return &pemPair
}

const privateKeyField string = "privateKey"
const certificateField string = "certificate"

// Reads the pair of TLS certificate and key from the specified secret.
func (c *Client) ReadTlsPair(props tls.TlsCertificateProps) *tls.TlsPemPair {
	name := generateSecretName(props)
	unstrSecret, err := c.GetResource(Secrets, props.Namespace, name)
	if err != nil {
		c.logger.Printf("Unable to get secret %s/%s: %s", props.Namespace, name, err)
		return nil
	}
	secret, err := convertToSecret(unstrSecret)
	if err != nil {
		return nil
	}
	pemPair := tls.TlsPemPair{
		Certificate: secret.Data[certificateField],
		PrivateKey:  secret.Data[privateKeyField],
	}
	if len(pemPair.Certificate) == 0 {
		c.logger.Printf("TLS Certificate not found in secret %s/%s", props.Namespace, name)
		return nil
	}
	if len(pemPair.PrivateKey) == 0 {
		c.logger.Printf("TLS PrivateKey not found in secret %s/%s", props.Namespace, name)
		return nil
	}
	return &pemPair
}

// Writes the pair of TLS certificate and key to the specified secret.
// Updates existing secret or creates new one.
func (c *Client) WriteTlsPair(props tls.TlsCertificateProps, pemPair *tls.TlsPemPair) error {
	name := generateSecretName(props)
	_, err := c.GetResource(Secrets, props.Namespace, name)
	if err != nil {
		secret := &v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: props.Namespace,
			},
			Data: map[string][]byte{
				certificateField: pemPair.Certificate,
				privateKeyField:  pemPair.PrivateKey,
			},
		}

		_, err := c.CreateResource(Secrets, props.Namespace, secret)
		if err == nil {
			c.logger.Printf("Secret %s is created", name)
		}
		return err
	}
	secret := v1.Secret{}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data[certificateField] = pemPair.Certificate
	secret.Data[privateKeyField] = pemPair.PrivateKey

	_, err = c.UpdateResource(Secrets, props.Namespace, secret)
	if err != nil {
		return err
	}
	c.logger.Printf("Secret %s is updated", name)
	return nil
}

func generateSecretName(props tls.TlsCertificateProps) string {
	return tls.GenerateInClusterServiceName(props) + ".kyverno-tls-pair"
}
