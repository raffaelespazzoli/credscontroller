package com.raffa;

import java.io.File;
import java.io.IOException;
import java.util.HashMap;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Lazy;
import org.springframework.core.env.ConfigurableEnvironment;
import org.springframework.core.env.MapPropertySource;
import org.springframework.core.env.MutablePropertySources;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;

@Configuration
@Lazy(false)
public class InitConfiguration {

	private Log log = LogFactory.getLog(InitConfiguration.class);

	@Value("${vault.token.file:/var/run/secrets/vaultproject.io}")
	String secretFile;

	@Autowired
	private ConfigurableEnvironment env;

	@Bean
	 public MapPropertySource secretPropertySource(String vaultToken) throws JsonParseException, JsonMappingException, IOException {
	    ObjectMapper mapper = new ObjectMapper(); 
	    File from = new File(secretFile); 
	    TypeReference<HashMap<String,Object>> typeRef = new TypeReference<HashMap<String,Object>>() {};
	    HashMap<String,Object> o = mapper.readValue(from, typeRef); 
	    MapPropertySource vaultSecretPropertySource= new MapPropertySource("vaultSecretPropertySource", o);
	    MutablePropertySources sources = env.getPropertySources();
	    sources.addFirst(vaultSecretPropertySource );
	    return vaultSecretPropertySource;
	 }
}