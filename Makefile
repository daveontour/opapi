DESTDIR := C:\Users\dave_\Desktop\opapi

build:
	-C:\Users\dave_\go/bin/golangci-lint.exe  run
	go build -o $(DESTDIR)\opapi.exe .\opapi\main
	go build -o $(DESTDIR)\opapiseeder.exe .\opapiseeder\main
	go build -o $(DESTDIR)\webhookclient.exe .\webhookclient\main
	copy *.json $(DESTDIR)
	-mkdir $(DESTDIR)\testfiles
	copy .\testfiles $(DESTDIR)\testfiles
	copy help.html $(DESTDIR)\help.html
	copy adminhelp.htm $(DESTDIR)\adminhelp.htm
	-mkdir $(DESTDIR)\adminhelp_files
	copy .\adminhelp_files $(DESTDIR)\adminhelp_files

clean:
	del $(DESTDIR)\opapi.exe
	del $(DESTDIR)\opapiseeder.exe
	del $(DESTDIR)\webhookclient.exe

