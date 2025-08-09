<template>
  <v-card class="about-card mx-auto" max-width="900">
    <v-card-title class="text-h4 text-center pa-6">
      <v-icon left color="primary" size="large">mdi-information</v-icon>
      Sobre o Projeto
    </v-card-title>
    
    <v-card-text>
      <v-container>
        <v-row align="center" justify="center">
          <!-- Profile Section -->
          <v-col cols="12" md="4" class="text-center">
            <v-avatar size="150" class="mb-4 profile-image">
              <v-img
                :src="profileImage"
                alt="Rafael Dantas"
              >
                <template v-slot:placeholder>
                  <v-skeleton-loader type="avatar"></v-skeleton-loader>
                </template>
              </v-img>
            </v-avatar>
                        
            <div>
              <h2 class="text-h5 mb-1">{{ developerName }}</h2>
              <h4 class="mb-1 text-subtitle-2 font-weight-thin">Autor do projeto</h4>
            </div>
            
            <v-btn
              :href="linkedinUrl"
              target="_blank"
              color="primary"
              variant="outlined"
              prepend-icon="mdi-linkedin"
              class="mb-4"
            >
              Conectar no LinkedIn
            </v-btn>
          </v-col>
          
          <!-- Animated Description Section -->
          <v-col cols="12" md="8">
            <v-card variant="tonal" class="pa-4">
              <div class="animated-text-container">
                <h3 class="text-h6 mb-3">
                  <v-icon color="primary" class="mr-2">mdi-lightbulb</v-icon>
                  Por que LocalStack Web?
                </h3>
                
                <div>
                  <p 
                    v-for="(paragraph, index) in paragraphs" 
                    :key="index"
                    class="mb-3"
                  >
                    {{ paragraph }}
                  </p>
                </div>
                
                <v-divider class="my-4"></v-divider>
                
                <div class="collaboration-section">
                  <h4 class="text-h6 mb-2">
                    <v-icon color="success" class="mr-2">mdi-handshake</v-icon>
                    Colaboração Bem-Vinda!
                  </h4>
                  <p class="mb-3">
                    Este projeto é open source e aceita contribuições da comunidade. 
                    Sinta-se à vontade para reportar bugs, sugerir melhorias ou contribuir com código.
                  </p>
                  
                  <v-row>
                    <v-col cols="12" sm="6">
                      <v-btn
                        color="success"
                        variant="outlined"
                        prepend-icon="mdi-github"
                        block
                        @click="openGitHub"
                      >
                        Ver no GitHub
                      </v-btn>
                    </v-col>
                    <v-col cols="12" sm="6">
                      <v-btn
                        color="info"
                        variant="outlined"
                        prepend-icon="mdi-bug"
                        block
                        @click="reportIssue"
                      >
                        Reportar Issue
                      </v-btn>
                    </v-col>
                  </v-row>
                </div>
              </div>
            </v-card>
          </v-col>
        </v-row>
        
        <!-- Technical Stack Section -->
        <v-row class="mt-6">
          <v-col cols="12">
            <v-card variant="tonal" class="pa-4">
              <h3 class="text-h6 mb-3">
                <v-icon color="primary" class="mr-2">mdi-tools</v-icon>
                Stack Tecnológico
              </h3>
              
              <v-row>
                <v-col 
                  v-for="tech in techStack" 
                  :key="tech.name"
                  cols="6" 
                  sm="4" 
                  md="3"
                >
                  <v-chip
                    :color="tech.color"
                    variant="elevated"
                    class="ma-1"
                    :prepend-icon="tech.icon"
                  >
                    {{ tech.name }}
                  </v-chip>
                </v-col>
              </v-row>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  name: 'AboutComponent',
  data() {
    return {
      developerName: 'Rafael Dantas',
      linkedinUrl: 'https://www.linkedin.com/in/rafael-dantas-a0ab7248/',
      profileImage: 'https://media.licdn.com/dms/image/v2/D4D03AQFyjuLTISyHxg/profile-displayphoto-shrink_800_800/B4DZdRvik4HwAg-/0/1749423107421?e=1756339200&v=beta&t=RT1QDbxQsPkTiUgGU7OCH0mftx7FoDpcT6baYghl12g',
      paragraphs: [
        'O LocalStack Web nasceu da necessidade de ter uma interface visual intuitiva para gerenciar serviços AWS localmente durante o desenvolvimento.',
        'Com esta ferramenta, desenvolvedores podem visualizar e interagir com serviços como S3, DynamoDB, Lambda e outros de forma mais eficiente.',
        'O projeto utiliza Vue.js 3 com Vuetify para criar uma experiência moderna e responsiva, facilitando o trabalho com LocalStack.',
        'A interface permite monitorar recursos, executar operações e debugar aplicações sem a necessidade de usar apenas a linha de comando.'
      ],
      techStack: [
        { name: 'Vue.js 3', color: 'success', icon: 'mdi-vuejs' },
        { name: 'Vuetify', color: 'indigo', icon: 'mdi-material-ui' },
        { name: 'LocalStack', color: 'orange', icon: 'mdi-aws' },
        { name: 'AWS SDK', color: 'warning', icon: 'mdi-aws' },
        { name: 'Vite', color: 'purple', icon: 'mdi-lightning-bolt' },
        { name: 'Docker', color: 'blue', icon: 'mdi-docker' },
        { name: 'JavaScript', color: 'yellow', icon: 'mdi-language-javascript' },
        { name: 'Node.js', color: 'green', icon: 'mdi-nodejs' }
      ]
    }
  },
  methods: {    
    openGitHub() {
      window.open('https://github.com/dantasrafael/localstack-web', '_blank')
    },
    
    reportIssue() {
      window.open('https://github.com/dantasrafael/localstack-web/issues', '_blank')
    }
  }
}
</script>

<style scoped>
.about-card {
  margin-top: 20px;
  border-radius: 16px;
}

.collaboration-section {
  background: rgba(var(--v-theme-surface), 0.05);
  border-radius: 12px;
  padding: 16px;
}

.v-chip {
  margin: 2px;
}

.profile-image {
  border: 3px solid rgb(var(--v-theme-primary));
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}
</style>
