package main

func main() {
	SeedData("locations", `INSERT INTO public.locations (name, active, company_id, address) VALUES ('admin_location', true, 1, 'admin_address');`)
}