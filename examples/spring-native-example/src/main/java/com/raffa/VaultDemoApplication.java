package com.raffa;

import javax.annotation.PostConstruct;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.vault.annotation.VaultPropertySource;


@SpringBootApplication
//@VaultPropertySource("secret/spring-example")
public class VaultDemoApplication {
	
	private Log log = LogFactory.getLog(VaultDemoApplication.class);

	public static void main(String[] args) {
		SpringApplication.run(VaultDemoApplication.class, args);
	}


    @Value("${password}")
    String password;
    
    @PostConstruct
    private void postConstruct() {
    	log.info("My password is: " + password);
    }
    
}
