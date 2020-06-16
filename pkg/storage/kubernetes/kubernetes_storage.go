package kubernetes

import (
	"encoding/json"
	"fmt"
	"github.com/nicolaferraro/connect/pkg/storage"
	"github.com/nicolaferraro/connect/pkg/token"
	"github.com/nicolaferraro/connect/pkg/util/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	labelTokenKey    = "connect.container-tools.io/token"
	labelProviderKey = "connect.container-tools.io/provider"
	labelGroupKey    = "connect.container-tools.io/group"
	connectInfoKey   = "connect-info"
	tokenKey         = "token"
)

type kubernetesTokenStorage struct {
	client    *v1.CoreV1Client
	namespace string
}

func NewKubernetesTokenStorage(namespace string) (storage.TokenStorage, error) {
	config, err := kubernetes.GetKubernetesConfig()
	if err != nil {
		return nil, err
	}
	client, err := v1.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &kubernetesTokenStorage{
		client:    client,
		namespace: namespace,
	}, nil
}

func (k *kubernetesTokenStorage) List() ([]string, error) {
	options := metav1.ListOptions{
		LabelSelector: labelTokenKey,
	}
	secrets, err := k.client.Secrets(k.namespace).List(options)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(secrets.Items))
	for _, s := range secrets.Items {
		names = append(names, s.Name)
	}
	return names, nil
}

func (k *kubernetesTokenStorage) Get(name string) (*token.Token, error) {
	secret, err := k.client.Secrets(k.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	// decode
	var tk *token.Token
	data := secret.Data[connectInfoKey]
	if len(data) == 0 {
		return nil, fmt.Errorf("no %q key found on secret %q", connectInfoKey, name)
	}
	if err := json.Unmarshal(data, &tk); err != nil {
		return nil, err
	}
	return tk, nil
}

func (k *kubernetesTokenStorage) Save(name string, tk *token.Token) error {
	data, err := json.Marshal(tk)
	if err != nil {
		return err
	}

	existing := true
	secret, err := k.client.Secrets(k.namespace).Get(name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		existing = false
		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: k.namespace,
				Name:      name,
			},
		}
	} else if err != nil {
		return err
	}

	if secret.Labels == nil {
		secret.Labels = make(map[string]string)
	}
	secret.Labels[labelTokenKey] = "true"
	secret.Labels[labelGroupKey] = tk.Provider.Group
	secret.Labels[labelProviderKey] = tk.Provider.ID

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data[connectInfoKey] = data
	secret.Data[tokenKey] = []byte(tk.GetAccessToken())

	if existing {
		_, err := k.client.Secrets(k.namespace).Update(secret)
		return err
	}

	_, err = k.client.Secrets(k.namespace).Create(secret)
	return err
}
