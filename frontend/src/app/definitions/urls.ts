export const apiURL = process.env.NODE_ENV === "production"
	? process.env.NEXT_PUBLIC_API_URL
	: "http://localhost:8080/api/steam/";

export const backendURL = process.env.NODE_ENV === "production"
	? process.env.NEXT_PUBLIC_BACKEND_URL
	: "http://localhost:8080/v1/"
