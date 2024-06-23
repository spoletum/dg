postgres "config" {
	connection = "pgsql://${HCLOUD_TOKEN}:${HCLOUD_DNS_TOKEN}@localhost:5432/foobar"

	table "example" {
		query = "SELECT * FROM example WHERE LAST_UPDATE BETWEEN ${upper(HCLOUD_DNS_TOKEN)} AND ${HCLOUD_TOKEN}"
	}
}