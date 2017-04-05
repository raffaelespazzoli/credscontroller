package com.raffa;

import java.io.File;
import java.io.IOException;
import java.util.HashMap;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cloud.bootstrap.config.PropertySourceLocator;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Lazy;
import org.springframework.core.env.Environment;
import org.springframework.core.env.MapPropertySource;
import org.springframework.core.env.PropertySource;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;

@Configuration
//@AutoConfigureBefore({ SpringLegacyExampleApplication.class })
@Lazy(false)
public class InitConfiguration implements PropertySourceLocator {

	private Log log = LogFactory.getLog(InitConfiguration.class);

	@Value("${secret.file:/var/run/secrets/vaultproject.io/secret.json}")
	String secretFile;

//	@Autowired
//	private ConfigurableEnvironment env;
//
//	 public MapPropertySource secretPropertySource() throws JsonParseException, JsonMappingException, IOException {
//	    ObjectMapper mapper = new ObjectMapper(); 
//	    File from = new File(secretFile); 
//	    TypeReference<HashMap<String,Object>> typeRef = new TypeReference<HashMap<String,Object>>() {};
//	    HashMap<String,Object> o = mapper.readValue(from, typeRef); 
//	    MapPropertySource secretPropertySource= new MapPropertySource("secretPropertySource", o);
//	    MutablePropertySources sources = env.getPropertySources();
//	    sources.addFirst(secretPropertySource );
//	    log.debug("added secretPropertySource: "+ o);
//	    return secretPropertySource;
//	 }
	 
	    @Override
	    public PropertySource<?> locate(Environment environment) {
	    	ObjectMapper mapper = new ObjectMapper(); 
		    File from = new File(secretFile); 
		    TypeReference<HashMap<String,Object>> typeRef = new TypeReference<HashMap<String,Object>>() {};
		    HashMap<String, Object> o=null;
			try {
				o = mapper.readValue(from, typeRef);
			} catch (IOException e) {
				// TODO Auto-generated catch block
				log.error(e);
			}
	        return new MapPropertySource("secretPropertySource", o);
	    }	 
	
//	@Override
//	public void afterPropertiesSet() throws Exception {
//		log.debug("After Properties Set");
//		secretPropertySource();
//	}
}