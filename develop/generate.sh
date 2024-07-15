
openssl req -newkey rsa:2048 \
  -new -nodes -x509 \
  -days 3650 \
  -out ${1}/cert.pem \
  -keyout ${1}/key.pem \
  -subj "/C=AR/ST=BuenosAires/L=CABA View/O=Balancify/OU=Balancify DEV/CN=balancify"

openssl req -newkey rsa:2048 \
  -new -nodes -x509 \
  -days 3650 \
  -out ${1}/client-cert.pem \
  -keyout ${1}/client-key.pem \
  -subj "/C=AR/ST=BuenosAires/L=CABA View/O=bff/OU=bff DEV/CN=bff"