package com.raffa;

import java.io.BufferedInputStream;
import java.io.ByteArrayOutputStream;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.Collection;
import java.util.Iterator;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.AutoConfigureBefore;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.cloud.vault.config.VaultBootstrapConfiguration;
import org.springframework.cloud.vault.config.VaultProperties;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Lazy;
import org.springframework.core.env.ConfigurableEnvironment;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.vault.support.VaultTokenResponse;

import com.fasterxml.jackson.databind.ObjectMapper;

@Configuration
@EnableConfigurationProperties({ VaultProperties.class })
@AutoConfigureBefore({ VaultBootstrapConfiguration.class })
@Lazy(false)
public class VaultInitConfiguration implements InitializingBean {

	private Log log = LogFactory.getLog(VaultInitConfiguration.class);

	@Value("${vault.token.file:/var/run/secrets/vaultproject.io}")
	String secretFile;

	@Value("${vault.cacert:/var/run/secrets/kubernetes.io/serviceaccount/service-ca.crt}")
	String vaultCaFile;

	// @Value("${vault.truststore:/var/run/secrets/vault.jks}")
	// String vaultTrustStore;

	@Value("${vault.truststore.password:changeit}")
	String truststorePassword;

	VaultProperties vaultProperties;

	public String vaultToken() throws FileNotFoundException, IOException {
		try {
			log.debug("My vault secretFile is: " + secretFile);

			byte[] jsonData = Files.readAllBytes(Paths.get(secretFile));
			// create ObjectMapper instance
			ObjectMapper objectMapper = new ObjectMapper();

			// convert json string to object
			VaultTokenResponse token = objectMapper.readValue(jsonData, VaultTokenResponse.class);

			return token.getToken().getToken();
		} finally {
			Files.deleteIfExists(Paths.get(secretFile));
		}
	}

	@Autowired
	private ConfigurableEnvironment env;

	// @Bean
	// public MapPropertySource vaultSecretPropertySource(String vaultToken) {
	// Map<String, Object> vaultPropertySourceMap=new HashMap<String,Object>();
	// vaultPropertySourceMap.put("spring.cloud.vault.token", vaultToken);
	// MapPropertySource vaultSecretPropertySource= new
	// MapPropertySource("vaultSecretPropertySource", vaultPropertySourceMap);
	// MutablePropertySources sources = env.getPropertySources();
	// sources.addFirst(vaultSecretPropertySource );
	// System.out.println("spring.cloud.vault.token property added with value:
	// "+vaultToken);
	// return vaultSecretPropertySource;
	// }

	public ByteArrayResource inMemoryVaultTrustStore()
			throws CertificateException, KeyStoreException, NoSuchAlgorithmException, IOException {
		FileInputStream fis = new FileInputStream(vaultCaFile);
		Collection<X509Certificate> cca = (Collection<X509Certificate>) CertificateFactory.getInstance("X.509")
				.generateCertificates(new BufferedInputStream(fis));
		KeyStore ks = KeyStore.getInstance(KeyStore.getDefaultType());
		ks.load(null, null);
		Iterator<X509Certificate> i = cca.iterator();
		int ii = 0;
		while (i.hasNext()) {
			ks.setCertificateEntry(String.valueOf(ii++), i.next());
		}
		ByteArrayOutputStream baos = new ByteArrayOutputStream();
		// FileOutputStream fos=new FileOutputStream(vaultTrustStore);
		ks.store(baos, truststorePassword.toCharArray());
		// ks.store(fos, truststorePassword.toCharArray());
		log.debug("In memory truststore created: " + ks);
		return new ByteArrayResource(baos.toByteArray());
	}

	public VaultInitConfiguration(VaultProperties vaultProperties) {
		this.vaultProperties = vaultProperties;
	}

	@Override
	public void afterPropertiesSet() throws Exception {
		log.debug("After Properties Set");
		if (vaultProperties.getSsl() != null && vaultProperties.getSsl().getTrustStore() == null) {
			vaultProperties.getSsl().setTrustStore(inMemoryVaultTrustStore());
		}
		vaultProperties.setToken(vaultToken());
	}
}