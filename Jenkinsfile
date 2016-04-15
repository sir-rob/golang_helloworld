node('docker') {
  	checkout scm

  	docker.withRegistry('https://quay.io', 'quay-registry') {
  		def helloworld = docker.build ("aleks_saul/hello_world:${env.BUILD_TAG}", ".")
  		helloworld.push()
    }  	
}