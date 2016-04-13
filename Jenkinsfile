node('docker') {
  	checkout scm

	def helloworld = docker.build "aleks_saul/hello_world:test"

  	docker.withRegistry('https://quay.io', 'aleks_saul+jenkins') {
       	helloworld.push()
       	helloworld.push '${env.BUILD_TAG}'
    }  	
}
