package com.raffa;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@RestController
@PropertySource("file://${secret.file}")
public class SpringLegacyExampleApplication {

	public static void main(String[] args) {
		SpringApplication.run(SpringLegacyExampleApplication.class, args);
	}
	
	@Value("${password}")
	String password;

	@RequestMapping("/secret")
	public String secret() {
		return "my secret is" + password;
	}	
}
