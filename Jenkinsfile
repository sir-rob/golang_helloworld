node {
	try {
		git '…' 

		def helloworld = docker.build "alekssaul/helloworld:${env.BUILD_TAG}"

  		docker.withRegistry('https://quay.io', 'aleks_saul+jenkins') {
        	helloworld.push()
        	helloworld.push '${env.BUILD_TAG}'
      }  
	} catch (e) {
    throw e
	}  
}
