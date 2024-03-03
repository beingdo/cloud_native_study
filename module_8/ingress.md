openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=chris.com/O=chris" -addext "subjectAltName = DNS:chris.com"

kubectl create secret tls chris-tls --cert=./tls.crt --key=./tls.key

