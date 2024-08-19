package controllers

import (
	"fmt"

	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

// connectToHedera establishes a connection to the Hedera network using the provided account ID and private key.
// It returns a configured Hedera client and an error if any part of the connection process fails.
func connectHedera() (*hedera.Client, error) {
	accountID, err := hedera.AccountIDFromString(os.Getenv("HEDERA_ACCOUNT_ID"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse account ID: %v", err)
	}

	privateKey, err := hedera.PrivateKeyFromString(os.Getenv("HEDERA_PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	client := hedera.ClientForTestnet() // Use ClientForMainnet() for the mainnet
	client.SetOperator(accountID, privateKey)

	return client, nil
}

func SaveDocumentDetailsOnBlock(client *hedera.Client, contractID hedera.ContractID, documentHash string, customerID string, state string, version string) (string, error) {
	// Define the function name you want to call on the smart contract
	functionName := "tokenizeDocument"

	// Create a new contract call query
	contractCallQuery := hedera.NewContractExecuteTransaction().
		SetContractID(contractID).
		SetGas(100000). // Set gas limit for the transaction
		SetFunction(
			functionName, // The name of the function in the smart contract
			hedera.NewContractFunctionParameters().
				AddString(documentHash). // Add document hash as a parameter
				AddString(customerID).   // Add customer ID as a parameter
				AddString(state).        // Add state as a parameter
				AddString(version),      // Add version as a parameter
		)

	// Execute the transaction
	txResponse, err := contractCallQuery.Execute(client)
	if err != nil {
		return "", fmt.Errorf("failed to execute contract call: %v", err)
	}

	// Fetch the receipt to confirm the transaction succeeded
	receipt, err := txResponse.GetReceipt(client)
	if err != nil {
		return "", fmt.Errorf("failed to get receipt: %v", err)
	}

	// Verify the transaction status
	if receipt.Status != hedera.StatusSuccess {
		return "", fmt.Errorf("transaction failed with status: %v", receipt.Status)
	}
	value := receipt.Status.String()
	// You can also retrieve the function's return value (if any)
	//record, err := txResponse.GetRecord(client)
	//if err != nil {
	//	return "", fmt.Errorf("failed to get record: %v", err)
	//}

	// Assume the contract function returns a single string
	//returnValue := record.ContractFunctionResult().GetString(0)
	// value :=record.GetContractExecuteResult()

	return value, nil
}
