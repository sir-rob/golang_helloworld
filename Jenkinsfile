node('docker') {
  	checkout scm
  	
  	stage 'Build Docker Canary Image'  	
  	def helloworld = docker.build ("aleks_saul/hello_world:canary", ".")
  	echo("${env.CVS_BRANCH}")
  	
  	# MY_SECRET_USER injected from username and password binding
# The username and password can be conjoined or separated
	sh("set +x")
	sh ("echo  $quay-registry_USER $quay-registry_PASSWORD")

  	docker.withRegistry('https://quay.io/v1', 'quay-registry') {
  		stage 'Push Canary Image to Quay'  		
  		helloworld.push()
    }


}