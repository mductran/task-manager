# golang makefile based on https://golangdocs.com/makefiles-golang
BINARY_NAME=main.exe
 
build:
    go build -o ${BINARY_NAME} main.go

    # This runs signtool with a cert in your profile store instead of a *.pfx file, to avoid needing to store a password in the makefile or environment variable: https://stackoverflow.com/questions/26998439/signtool-with-certificate-stored-in-local-computer
    signtool sign /sm /s My /n <certificateSubjectName> /t http://timestamp.digicert.com ${BINARY_NAME}
 
run:
    go build -o ${BINARY_NAME} main.go
    ./${BINARY_NAME}
 
clean:
    go clean
    rm ${BINARY_NAME}