package main

import (
	"fmt"
	"github.com/nicolaferraro/connect/pkg/storage/kubernetes"
	"os"
	"time"
)

func main() {
	fmt.Println("Connect agent started")

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		fmt.Println("NAMESPACE environment variable not provided")
		os.Exit(1)
	}

	store, err := kubernetes.NewKubernetesTokenStorage(namespace)
	if err != nil {
		panic(err)
	}

	first := true
	for {
		if !first {
			time.Sleep(30 * time.Second)
		}
		first = false

		lst, err := store.List()
		if err != nil {
			fmt.Printf("ERROR: cannot list stored tokens: %v\n", err)
			continue
		}
		if len(lst) == 0 {
			fmt.Printf("No tokens found in namespace %q\n", namespace)
			continue
		}
		for _, tkName := range lst {
			tk, err := store.Get(tkName)
			if err != nil {
				fmt.Printf("ERROR: cannot get tokens %q: %v\n", tkName, err)
				continue
			}
			newToken, err := tk.Refresh()
			if err != nil {
				fmt.Printf("ERROR: cannot refresh token %q: %v\n", tkName, err)
				continue
			}
			if newToken.GetAccessToken() != tk.GetAccessToken() {
				// Store the new token
				err = store.Save(tkName, newToken)
				if err != nil {
					fmt.Printf("ERROR: cannot save the new credentials for token %q: %v\n", tkName, err)
					continue
				}
				fmt.Printf("Token %q has been refreshed. New expiry date is %v", tkName, newToken.GetExpiry())
			} else {
				fmt.Printf("No need to refresh token %q until %v\n", tkName, newToken.GetExpiry())
			}
		}
	}

}
