node('docker') {
  	checkout scm

	def helloworld = docker.build ("aleks_saul/hello_world:${env.BUILD_TAG}", "aleks_saul/hello_world"

  	docker.withRegistry('https://quay.io', 'aleks_saul+jenkins') {
  		helloworld.push()
       	helloworld.push '${env.BUILD_TAG}'
    }  	
}