DESTDIR := C:\Users\dave_\Desktop\opapi


build:
#	-C:\Users\dave_\go/bin/golangci-lint.exe  run
	go build -o $(DESTDIR)\opapi.exe .\opapi\main
	go build -o $(DESTDIR)\opapiseeder.exe .\opapiseeder\main
	go build -o $(DESTDIR)\webhookclient.exe .\webhookclient\main
	go build -o $(DESTDIR)\rabbitclient.exe .\rabbitclient\main
	go build -o $(DESTDIR)\perftestclient.exe .\perftestclient\main
	copy *.json $(DESTDIR)
	-mkdir $(DESTDIR)\testfiles
	copy .\testfiles $(DESTDIR)\testfiles


clean:
	del $(DESTDIR)\opapi.exe
	del $(DESTDIR)\opapiseeder.exe
	del $(DESTDIR)\webhookclient.exe
	del $(DESTDIR)\rabbitclient.exe
	del $(DESTDIR)\perftestclient.exe
