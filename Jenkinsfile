node('docker') {
  	checkout scm

	docker.withRegistry('https://quay.io', 'aleks_saul+jenkins') {
		def helloworld = docker.build ("aleks_saul/hello_world:${env.BUILD_TAG}", "aleks_saul/hello_world"
  		helloworld.push()
       	helloworld.push '${env.BUILD_TAG}'
    }  	
}