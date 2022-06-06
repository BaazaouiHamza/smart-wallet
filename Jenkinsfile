@Library('jenkins-shared-libraries')

def project_url = "https://git.digitus.me/pfe/smart-wallet.git"

properties(
    [   
        gitLabConnection('jenkins'),
        buildDiscarder(
            logRotator(
                daysToKeepStr: '60',
                numToKeepStr: '200'
            )
        ),
        disableConcurrentBuilds()
    ]
)    
    
timestamps {
    node(){
        golangPipeline(project_url)
    }
}
