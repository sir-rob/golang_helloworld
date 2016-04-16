node('docker') {
  	checkout scm
  	
  	stage 'Build Docker Canary Image'  	
  	def helloworld = docker.build ("aleks_saul/hello_world:canary", ".")
  	
  	docker.withRegistry('https://quay.io/v1', 'quay-registry') {
  		stage 'Push Canary Image to Quay'  		
  		helloworld.push()
    }


}
