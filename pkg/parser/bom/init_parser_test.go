package bom

// func TestInitParsers(t *testing.T) {

// 	// Test case 1: Image argument provided
// 	arg1 := model.NewArguments()
// 	arg1.Image = stringPtr("alpine")

// 	InitParsers(*arg1)
// 	if !assert.DirExists(t, *Target) {
// 		t.Errorf("Target Directory for image %s was not set correctly '%s'", *arg1.Image, *Target)
// 	}

// 	// Test case 2: Dir argument provided
// 	arg2 := model.NewArguments()
// 	arg2.Dir = stringPtr(".")

// 	InitParsers(*arg2)
// 	if !assert.DirExists(t, *Target) {
// 		t.Errorf("Target Directory for %s was not set correctly '%s'", *arg2.Dir, *Target)
// 	}

// 	// Test case 3: Tar argument provided
// 	tarFile := docker.SaveImageToTar(stringPtr("alpine"))
// 	arg3 := model.NewArguments()
// 	arg3.Tar = stringPtr(tarFile.Name())

// 	InitParsers(*arg3)
// 	if !assert.DirExists(t, *Target) {
// 		t.Errorf("Target Directory for tar file %s was not set correctly '%s'", *arg3.Tar, *Target)
// 	}
// }

// // Helper function to create a string pointer
// func stringPtr(s string) *string {
// 	return &s
// }
