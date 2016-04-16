node('docker') {
  	checkout scm
  	  	 		
  	def containertag = "canary" 
  	
	if (env.BRANCH_NAME == 'master') {
    	containertag = "production" 
	}

	stage 'Build Docker Image'  	
  	def helloworld = docker.build ("aleks_saul/hello_world:$containertag", ".")
	
	stage 'Push docker Image to Quay' 
  	docker.withRegistry('https://quay.io/v1', 'quay-registry') {
  		 		
  		helloworld.push()
    }


}
