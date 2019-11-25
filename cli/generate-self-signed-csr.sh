#ssh-keygen -t rsa -b 2048 -f pk.pem -q
openssl req \
       -newkey rsa:2048 -sha256 -nodes -keyout dip.key \
       -x509 -days 365 -out dip.crt \
       -subj '/CN=and/C=DE'
