package util

import "testing"

func TestZipCompressDir(t *testing.T) {

}

func TestZipDeCompress(t *testing.T) {
	zip := "/Users/king/Downloads/Specs-master.zip"
	dst := "/Users/king/Downloads/SpecsMaster"
	err := ZipDeCompress(zip, dst)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTarGzDir(t *testing.T) {
	dir := "/Users/king/Downloads/cocoapods-bin-master"
	tarfile := "/Users/king/Downloads/cocoapods-bin-master.tar.gz"
	err := TarGzDir(dir, tarfile)
	if err != nil {
		t.Fatal(err)
	}
}
