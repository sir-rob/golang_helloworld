node('docker') {
  	checkout scm
  	stage 'Build Docker Image'
  	sh "echo Building the Docker Image"
  	def helloworld = docker.build ("aleks_saul/hello_world:${env.BUILD_TAG}", ".")
  	sh "echo Built the Docker Image"

  	docker.withRegistry('https://quay.io', "quay-registry") {
  		sh "echo Building the quay-registry"
  		stage 'Push Image'
  		helloworld.push()
    }  	
}