package main

import (
	"errors"
	"os"

	vault "github.com/hashicorp/vault/api"
)

type Output interface{}

type Executer interface {
	Execute() (Output, error)
}

type Action interface {
	Executer
	Parse(context map[string]interface{}) error // TODO: remove from here?
}

// -------------------------------------------------------------

type ActionCreator func() Action

func NewWriteFileAction() Action {
	return &WriteFileAction{}
}

type WriteFileAction struct {
	File  string `json:"file"`
	Value string `json:"value"`
}

func (w *WriteFileAction) Execute() (Output, error) {
	f, err := os.Create(w.File)
	if err != nil {
		return nil, err
	}
	if _, err = f.WriteString(w.Value); err != nil {
		return nil, err
	}
	return nil, f.Sync()
}

func (w *WriteFileAction) Parse(context map[string]interface{}) error {
	parsedValue, err := parse(w.File, context)
	if err != nil {
		return err
	}
	w.File = parsedValue.(string)
	parsedValue, err = parse(w.Value, context)
	if err != nil {
		return err
	}
	w.Value = parsedValue.(string)
	return nil
}

func NewVaultLogicalAction() Action {
	return &VaultLogicalAction{}
}

type VaultLogicalAction struct {
	Token  string                 `json:"token,omitempty"`
	Method string                 `json:"method"`
	Path   string                 `json:"path"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

func (l *VaultLogicalAction) Execute() (Output, error) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.ReadEnvironment()
	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}
	switch l.Method {
	case "read":
		return vaultClient.Logical().Read(l.Path)
	case "write":
		return vaultClient.Logical().Write(l.Path, l.Data)
	}
	return nil, errors.New("Method not allowed")
}

func (l *VaultLogicalAction) Parse(context map[string]interface{}) error {
	parsedValue, err := parse(l.Path, context)
	if err != nil {
		return err
	}
	l.Path = parsedValue.(string)
	if l.Data != nil {
		parsedValue, err = parse(l.Data, context)
		if err != nil {
			return err
		}
		l.Data = parsedValue.(map[string]interface{})
	}
	return nil
}

func GetAvailableActions() map[string]ActionCreator {
	return map[string]ActionCreator{
		"writeFile":    NewWriteFileAction,
		"vaultLogical": NewVaultLogicalAction,
		// ... generate this map automatically
	}
}
