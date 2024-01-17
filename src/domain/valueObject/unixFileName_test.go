package valueObject

import "testing"

func TestUnixFileName(t *testing.T) {
	t.Run("ValidUnixFileName", func(t *testing.T) {
		validUnixFileNames := []string{
			"17795713_1253219528108045_4440713319482755723_n\\ \\(1\\).png",
			"hello.php",
			"hello_file.php",
			"hello\\$_file.php",
			"hello (file).php",
			"Imagem - Sem Título.jpg",
			"Imagem - Sem Título & BW.jpg",
			"Imagem - Sem Título # BW.jpg",
			"Imagem - Sem Título @ BW.jpg",
			"Clean Architecture A Craftsman's Guide to Software Structure and Design.pdf",
			".sudo_as_admin_successful",
			"WhatsApp Image 2018-06-22 at 18.05.08.jpeg",
		}
		for _, name := range validUnixFileNames {
			_, err := NewUnixFileName(name)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", name, err)
			}
		}
	})

	t.Run("InvalidUnixFileName", func(t *testing.T) {
		invalidUnixFileNames := []string{
			"",
			".",
			"..",
			"/",
			"\\",
			"file.php?blabla",
			"@<php52.sandbox.ntorga.com>.php",
			"../file.php",
			"hello10/info.php",
		}
		for _, name := range invalidUnixFileNames {
			_, err := NewUnixFileName(name)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", name)
			}
		}
	})
}
