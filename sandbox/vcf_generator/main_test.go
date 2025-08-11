package vcfgenerator

import (
	"fmt"
	"os"
	"testing"
)

// Test_GenerateVCF generates .vcf file with a list of numbers to add to the block list on iphone
func Test_GenerateVCF(t *testing.T) {
	prefix := "+38044390"

	for batch := range 10 {

		vcard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:Block List %s%d***\nN:Block %s%d***;Spam;;;\n", prefix, batch, prefix, batch)

		for i := range 1000 {
			vcard += fmt.Sprintf("TEL;TYPE=CELL:%s%d%03d\n", prefix, batch, i)
		}
		vcard += "END:VCARD\n"

		file, err := os.OpenFile(fmt.Sprintf("blocklist_%s_%d.vcf", prefix[1:], batch), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			t.Fatalf("failed to open file: %v", err)
		}

		if _, err := file.WriteString(vcard); err != nil {
			t.Fatalf("failed to write to file: %v", err)
		}
		file.Close()
	}
}
