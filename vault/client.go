package vault

import (
	"errors"
	"fmt"
	"io/ioutil"

	. "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

type Vault struct {
	Host           string
	Port           string
	Scheme         string
	Role           string
	Mount          string
	Namespace      string
	ServiceAccount string
}

type msiResponseJson struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}

var client *Client

func (v *Vault) Initialize() error {
	var err error
	var renew bool
	var token string

	//Default client
	config := DefaultConfig()
	client, err = NewClient(config)
	//Set the address
	err = client.SetAddress(fmt.Sprintf("%s://%s:%s", v.Scheme, v.Host, v.Port))
	if err != nil {
		return err
	}

	//Set the namespace if we have one
	if len(v.Namespace) > 0 {
		log.Println("Namespace: " + v.Namespace)
		client.SetNamespace(v.Namespace)
	}

	//Auth to Vault
	log.Println("Client authenticating to Vault")
	log.Println("Using kubernetes authentication")

	//Check Mount
	if len(v.Mount) == 0 {
		return errors.New("Auth mount not in config.")
	}
	log.Printf("Mount: auth/%s", v.Mount)

	//Check Role
	if len(v.Role) == 0 {
		return errors.New("K8s role not in config.")
	}
	log.Printf("Role: %s", v.Role)

	//Check SA
	if len(v.ServiceAccount) == 0 {
		return errors.New("K8s SA file not in config.")
	}
	log.Printf("SA: %s", v.ServiceAccount)

	//Get the JWT from POD
	jwt, err := ioutil.ReadFile(v.ServiceAccount)
	if err != nil {
		return err
	}

	//Auth with K8s vault
	data := map[string]interface{}{"jwt": string(jwt), "role": v.Role}
	secret, err := client.Logical().Write(fmt.Sprintf("auth/%s/login", v.Mount), data)
	if err != nil {
		return err
	}

	//Set client token
	token = secret.Auth.ClientToken
	client.SetToken(token)

	//See if the token we got is renewable
	log.Println("Looking up token")
	lookup, err := client.Auth().Token().LookupSelf()
	//If token is not valid so get out of here early
	if err != nil {
		return err
	}

	//Check renewable
	renew = lookup.Data["renewable"].(bool)
	if renew == true {
		go v.RenewToken()
	}

	return nil
}

func (v *Vault) GetSecret(path string) (Secret, error) {
	log.Printf("Getting secret: %s", path)
	secret, err := client.Logical().Read(path)
	if err != nil {
		return Secret{}, err
	}
	return *secret, nil
}

func (v *Vault) PutSecret(path string, secret map[string]interface{}) error {
	log.Printf("Getting secret: %s", path)
	_, err := client.Logical().Write(path, secret)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) GetCertificate(path string, cn string) (Secret, error) {
	log.Printf("Getting certificate: %s", path)
	certData := map[string]interface{}{
		"common_name": cn,
	}
	secret, err := client.Logical().Write(path, certData)
	if err != nil {
		return Secret{}, err
	}
	return *secret, nil
}

func (v *Vault) RenewToken() {
	//If it is let's renew it by creating the payload
	secret, err := client.Auth().Token().RenewSelf(0)
	if err != nil {
		log.Fatal(err)
	}

	//Create the object. TODO look at setting increment explicitly
	renewer, err := client.NewRenewer(&RenewerInput{
		Secret: secret,
		//Grace:  time.Duration(15 * time.Second),
		//Increment: 60,
	})

	//Check if we were able to create the renewer
	if err != nil {
		log.Fatal(err)
	}

	//Start the renewer
	log.Printf("Starting token lifecycle management for accessor: %s", secret.Auth.Accessor)
	go renewer.Renew()
	defer renewer.Stop()

	//Log it
	for {
		select {
		case err := <-renewer.DoneCh():
			if err != nil {
				log.Fatal(err)
			}
			//App will terminate after token cannot be renewed.
			log.Fatalf("Cannot renew token with accessor %s. App will terminate.", secret.Auth.Accessor)
		case renewal := <-renewer.RenewCh():
			log.Printf("Successfully renewed token accessor: %s", renewal.Secret.Auth.Accessor)
		}
	}
}

func (v *Vault) RenewSecret(secret Secret) error {
	renewer, err := client.NewRenewer(&RenewerInput{
		Secret: &secret,
		//Grace:  time.Duration(15 * time.Second),
	})

	//Check if we were able to create the renewer
	if err != nil {
		log.Fatal(err)
	}

	//Start the renewer
	log.Printf("Starting secret lifecycle management for lease: %s", secret.LeaseID)
	go renewer.Renew()
	defer renewer.Stop()

	//Log it
	for {
		select {
		case err := <-renewer.DoneCh():
			if err != nil {
				log.Fatal(err)
			}
			//Renewal is now past max TTL. Let app die reschedule it elsewhere. TODO: Allow for getting new creds here.
			log.Fatalf("Cannot renew %s. App will terminate.", secret.LeaseID)
		case renewal := <-renewer.RenewCh():
			log.Printf("Successfully renewed secret lease: %s", renewal.Secret.LeaseID)
		}
	}
}

func (v *Vault) Encrypt(path string, plaintext string) (string, error) {
	var ciphertext string

	data := map[string]interface{}{"plaintext": plaintext}
	secret, err := client.Logical().Write(path, data)
	if err != nil {
		return "", err
	}

	ciphertext = secret.Data["ciphertext"].(string)
	return ciphertext, nil
}

func (v *Vault) Decrypt(path string, ciphertext string) (string, error) {
	var plaintext string

	data := map[string]interface{}{"ciphertext": ciphertext}
	secret, err := client.Logical().Write(path, data)
	if err != nil {
		return "", err
	}

	plaintext = secret.Data["plaintext"].(string)
	return plaintext, nil
}

func (v *Vault) Close() {
	client.Auth().Token().RevokeSelf(client.Token())
}
