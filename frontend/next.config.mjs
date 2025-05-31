/** @type {import('next').NextConfig} */
const nextConfig = {
	images: {
		remotePatterns: [
			{
				protocol: "https",
				hostname: "avatars.steamstatic.com",
				port: "",
				pathname: "/**",
			},
			{
				protocol: "http",
				hostname: "media.steampowered.com",
				port: "",
				pathname: "/**",
			}
		]
	},
	output: "standalone"
};

export default nextConfig;