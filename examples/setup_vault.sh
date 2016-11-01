### Vault pki bootstrap

vault mount -path=rca pki
vault mount-tune -default-lease-ttl=43800h -max-lease-ttl=87600h rca
vault mount -path=ica pki
vault mount-tune -default-lease-ttl=35040h -max-lease-ttl=70080h ica

vault write rca/root/generate/internal common_name=rootca ttl=87600h key_bits=4096
vault write -field=csr ica/intermediate/generate/internal common_name=interca ttl=70080h > inter.csr
cat inter.csr | vault write -field=certificate rca/root/sign-intermediate csr=- use_csr_values=true > inter.cert
rm inter.csr
cat inter.cert | vault write ica/intermediate/set-signed certificate=-
rm inter.cert

####### role
vault write ica/roles/default ttl=8760h allow_any_name=true
