package com.raffa;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.PropertySource;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@SpringBootApplication
@RestController
//@PropertySource("file://${secret.file}")
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
