package gpg

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/franela/vault/executor"
	"github.com/franela/vault/vault"
)

var logger = &logWriter{}

var cmdExec executor.Executor

func init() {
	cmdExec = &executor.CMDExecutor{}

}

func getGPGHomeDir() []string {
	if len(os.Getenv("GNUPGHOME")) > 0 {
		return []string{"--homedir", os.Getenv("GNUPGHOME")}
	}
	return []string{}
}

type logWriter struct {
}

func (*logWriter) Write(input []byte) (n int, err error) {
	log.Printf("%s", input)
	return len(input), nil

}

func Decrypt(filePath string) (string, error) {
	decryptArgs := append(getGPGHomeDir(), "--decrypt", "--armor", "--batch", "--yes", filePath)

	cmd := exec.Command("gpg", decryptArgs...)
	cmd.Env = nil
	cmd.Stderr = logger

	out, err := cmdExec.Output(cmd)

	if err != nil {
		return "", err
	}

	return string(out), nil
}

func DecryptFile(outputFile, filePath string) error {

	decryptArgs := append(getGPGHomeDir(), "--decrypt", "--armor", "--batch", "--yes", "--output", outputFile, filePath)

	cmd := exec.Command("gpg", decryptArgs...)
	cmd.Env = nil
	cmd.Stderr = logger

	_, err := cmdExec.Output(cmd)

	if err != nil {
		return err
	}

	return nil
}

func Encrypt(filePath string, text string, recipients []vault.VaultRecipient) error {
	if err := os.MkdirAll(path.Dir(filePath), 0700); err != nil {
		return err
	}

	encryptArgs := append(getGPGHomeDir(), "--encrypt", "--armor", "--batch", "--yes", "--output", filePath)

	for _, recipient := range recipients {
		encryptArgs = append(encryptArgs, "--recipient")
		encryptArgs = append(encryptArgs, recipient.Fingerprint)
	}

	log.Printf("Running: gpg %s\n", strings.Join(encryptArgs, " "))
	cmd := exec.Command("gpg", encryptArgs...)
	cmd.Env = nil
	cmd.Stderr = logger
	cmd.Stdin = strings.NewReader(text)
	_, err := cmdExec.Output(cmd)

	if err != nil {
		return err
	}
	return nil
}

func EncryptFile(filePath string, sourceFile string, recipients []vault.VaultRecipient) error {

	if err := os.MkdirAll(path.Dir(filePath), 0700); err != nil {
		return err
	}

	encryptArgs := append(getGPGHomeDir(), "--encrypt", "--armor", "--batch", "--yes", "--output", filePath)

	for _, recipient := range recipients {
		encryptArgs = append(encryptArgs, "--recipient")
		encryptArgs = append(encryptArgs, recipient.Fingerprint)
	}

	encryptArgs = append(encryptArgs, sourceFile)

	cmd := exec.Command("gpg", encryptArgs...)
	cmd.Env = nil
	cmd.Stderr = logger
	_, err := cmdExec.Output(cmd)

	if err != nil {
		return err
	}
	return nil
}

func ReEncryptFile(src, dst string, recipients []vault.VaultRecipient) error {
	decryptArgs := append(getGPGHomeDir(), "--decrypt", "--armor", "--batch", "--yes", src)
	encryptArgs := append(getGPGHomeDir(), "--encrypt", "--armor", "--batch", "--yes", "--output", dst)

	for _, recipient := range recipients {
		encryptArgs = append(encryptArgs, "--recipient")
		encryptArgs = append(encryptArgs, recipient.Fingerprint)
	}

	decryptCmd := exec.Command("gpg", decryptArgs...)
	decryptCmd.Env = nil
	encryptCmd := exec.Command("gpg", encryptArgs...)
	encryptCmd.Env = nil

	srcFile, fileErr := os.Open(src)
	if fileErr != nil {
		return fileErr
	}

	stat, statErr := srcFile.Stat()
	if statErr != nil {
		return statErr
	}

	r, w := io.Pipe()
	bufferedStdout := bufio.NewWriterSize(w, int(stat.Size()))
	decryptCmd.Stdout = bufferedStdout
	decryptCmd.Stderr = logger

	encryptCmd.Stderr = logger
	encryptCmd.Stdin = r

	err1 := cmdExec.Run(decryptCmd)
	if err1 != nil {
		return err1
	}
	w.Close()

	err2 := cmdExec.Run(encryptCmd)
	if err2 != nil {
		return err2
	}

	return nil
}

func ReceiveKey(recipients []vault.VaultRecipient) error {
	recvArgs := append(getGPGHomeDir(), "--batch", "--yes", "--recv-keys")

	for _, recipient := range recipients {
		recvArgs = append(recvArgs, recipient.Fingerprint)
	}

	recvCmd := exec.Command("gpg", recvArgs...)
	recvCmd.Env = nil
	recvCmd.Stderr = logger

	return cmdExec.Run(recvCmd)
}

func ReceiveKeyFromKeyserver(recipients []vault.VaultRecipient, keyserver string) error {
	recvArgs := append(getGPGHomeDir(), "--batch", "--yes", "--keyserver", keyserver, "--recv-keys")

	for _, recipient := range recipients {
		recvArgs = append(recvArgs, recipient.Fingerprint)
	}

	recvCmd := exec.Command("gpg", recvArgs...)
	recvCmd.Env = nil
	recvCmd.Stderr = logger

	return cmdExec.Run(recvCmd)
}

func DeleteKey(recipient vault.VaultRecipient) error {
	recvArgs := append(getGPGHomeDir(), "--batch", "--yes", "--delete-key", recipient.Fingerprint)

	recvCmd := exec.Command("gpg", recvArgs...)
	recvCmd.Env = nil
	recvCmd.Stderr = logger
	return cmdExec.Run(recvCmd)

}

func GetSecretKeysCount() (int, error) {
	recvArgs := append(getGPGHomeDir(), "--list-secret-keys", "--with-colons")

	recvCmd := exec.Command("gpg", recvArgs...)
	recvCmd.Env = nil
	recvCmd.Stderr = logger
	output, err := cmdExec.Output(recvCmd)

	if err != nil {
		return 0, err
	}

	return strings.Count(output, "\n"), nil

}

func GetKeysWithFingerprints() (map[string]bool, error) {
	recvargs := append(getGPGHomeDir(), "--list-keys", "--with-colons", "--fingerprint")

	recvcmd := exec.Command("gpg", recvargs...)
	recvcmd.Env = nil
	recvcmd.Stderr = logger
	output, err := cmdExec.Output(recvcmd)

	var fingerprintMap = map[string]bool{}

	if err != nil {
		return fingerprintMap, err
	}

	lines := strings.Split(output, "\n")

	for index, line := range lines {
		if strings.HasPrefix(line, "fpr") {
			keyValidity := strings.Split(lines[index-1], ":")[1]
			isTrusty := keyValidity == "u" || keyValidity == "f"
			fingerprintMap[strings.Split(line, ":")[9]] = isTrusty
		}
	}

	return fingerprintMap, nil

}

func GetKeyringOwnerRecipient() (*vault.VaultRecipient, error) {
	recvargs := append(getGPGHomeDir(), "--list-secret-keys", "--with-colons", "--fingerprint")

	recvcmd := exec.Command("gpg", recvargs...)
	recvcmd.Env = nil
	recvcmd.Stderr = logger
	output, err := cmdExec.Output(recvcmd)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")

	if len(lines) < 2 {
		return nil, errors.New("Keyring owner not found")
	}

	recipient := &vault.VaultRecipient{}

	recipient.Name = strings.Split(lines[0], ":")[9]
	recipient.Fingerprint = strings.Split(lines[1], ":")[9]

	return recipient, nil

}

func SetExecutor(executor executor.Executor) {
	cmdExec = executor
}
