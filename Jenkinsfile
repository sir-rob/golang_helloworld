node('docker') {
  	checkout scm
  	stage 'Build Docker Image'
  	
  	def helloworld = docker.build ("aleks_saul/hello_world:${env.BUILD_TAG}", ".")
  	
  	docker.withRegistry('https://quay.io/v1', "quay-registry") {
  		stage 'Push Image'
  		helloworld.push()
    }  	
}