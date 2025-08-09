<template>
  <div>
    <v-row>
      <v-col cols="12" class="bg-secondary opacity-50">
        <div class="d-flex justify-space-between align-center">
          <div class="d-flex align-center">
            <v-icon class="mr-2">mdi-forum</v-icon>
            <h1 class="text-h5">SNS Topics</h1>
          </div>
          <v-btn color="primary" @click="createTopicDialog = true">
            <v-icon class="mr-2">mdi-plus</v-icon>
            Criar Tópico
          </v-btn>
        </div>
      </v-col>
    </v-row>

    <!-- Search and Filter -->
    <v-row>
      <v-col cols="12">
        <v-text-field
          v-model="searchQuery"
          label="Buscar tópicos"
          prepend-inner-icon="mdi-magnify"
          variant="outlined"
          clearable
        >
          <template v-slot:append>
            <v-btn @click="loadTopics" :loading="loading" variant="outlined">
              <v-icon class="mr-2">mdi-refresh</v-icon>
              Atualizar
            </v-btn>
          </template>      
      </v-text-field>
      </v-col>
    </v-row>

    <!-- Topics List -->
    <v-row>
      <v-col cols="12" md="6" lg="4" v-for="topic in filteredTopics" :key="topic.TopicArn">
        <v-card class="mb-4" elevation="2">
          <v-card-title class="d-flex justify-space-between align-center">
            <TitleNameWithTooltip :name="getTopicName(topic.TopicArn)" />
            <v-menu>
              <template v-slot:activator="{ props }">
                <v-btn icon="mdi-dots-vertical" v-bind="props" size="small"></v-btn>
              </template>
              <v-list>
                <v-list-item @click="openTopic(topic.TopicArn)">
                  <v-list-item-title>
                    <v-icon class="mr-2">mdi-forum</v-icon>
                    Detalhes
                  </v-list-item-title>
                </v-list-item>
                <v-list-item @click="publishToTopic(topic.TopicArn)">
                  <v-list-item-title>
                    <v-icon class="mr-2">mdi-send</v-icon>
                    Publicar Mensagem
                  </v-list-item-title>
                </v-list-item>
                <v-list-item @click="confirmDeleteTopic(topic.TopicArn)">
                  <v-list-item-title class="text-error">
                    <v-icon class="mr-2">mdi-delete</v-icon>
                    Excluir
                  </v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-card-title>
          <v-card-text>
            <v-chip size="small" class="mb-2">
              <v-icon start>mdi-identifier</v-icon>
              ARN
            </v-chip>
            <p class="text-caption text-break">{{ topic.TopicArn }}</p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Empty State -->
    <v-row v-if="!loading && topics.length === 0">
      <v-col cols="12" class="text-center">
        <v-icon size="120" color="grey-lighten-2">mdi-forum-outline</v-icon>
        <h2 class="text-h5 text-grey-lighten-1 mt-4">Nenhum tópico encontrado</h2>
        <p class="text-grey-lighten-1">Crie seu primeiro tópico SNS para começar</p>
        <v-btn color="primary" @click="createTopicDialog = true" class="mt-4">
          <v-icon class="mr-2">mdi-plus</v-icon>
          Criar Tópico
        </v-btn>
      </v-col>
    </v-row>

    <!-- Loading State -->
    <v-row v-if="loading">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate size="64" color="primary"></v-progress-circular>
        <p class="mt-4 text-grey-lighten-1">Carregando tópicos...</p>
      </v-col>
    </v-row>

    <!-- Create Topic Dialog -->
    <v-dialog v-model="createTopicDialog" max-width="500">
      <v-card>
        <v-card-title>
          <span class="text-h5">Criar Novo Tópico</span>
        </v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newTopicName"
            label="Nome do Tópico"
            placeholder="meu-topico"
            variant="outlined"
            :rules="[rules.required]"
            @keyup.enter="createTopic"
          ></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="createTopicDialog = false">Cancelar</v-btn>
          <v-btn 
            color="primary" 
            @click="createTopic"
            :loading="creating"
            :disabled="!newTopicName"
          >
            Criar
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Publish Message Dialog -->
    <v-dialog v-model="publishDialog" max-width="600">
      <v-card>
        <v-card-title>
          <span class="text-h5">Publicar Mensagem</span>
        </v-card-title>
        <v-card-text>
          <v-text-field
            v-model="messageSubject"
            label="Assunto"
            placeholder="Assunto da mensagem"
            variant="outlined"
            class="mb-4"
          ></v-text-field>
          <v-textarea
            v-model="messageContent"
            label="Mensagem"
            placeholder="Conteúdo da mensagem"
            variant="outlined"
            rows="5"
            :rules="[rules.required]"
          ></v-textarea>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="publishDialog = false">Cancelar</v-btn>
          <v-btn 
            color="primary" 
            @click="publishMessage"
            :loading="publishing"
            :disabled="!messageContent"
          >
            Publicar
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Topic Details Dialog -->
    <v-dialog v-model="topicDetailsDialog" max-width="800">
      <v-card v-if="selectedTopic">
        <v-card-title>
          <span class="text-h5">Detalhes do Tópico</span>
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <v-card variant="outlined">
                <v-card-text>
                  <p><strong>Nome:</strong> {{ getTopicName(selectedTopic.TopicArn) }}</p>
                  <p><strong>ARN:</strong> {{ selectedTopic.TopicArn }}</p>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
          
          <v-row class="mt-4">
            <v-col cols="12">
              <h3 class="text-h6 mb-4">Subscriptions</h3>
              <v-list v-if="subscriptions.length > 0">
                <v-list-item v-for="subscription in subscriptions" :key="subscription.SubscriptionArn">
                  <v-list-item-title>{{ subscription.Protocol }}</v-list-item-title>
                  <v-list-item-subtitle>{{ subscription.Endpoint }}</v-list-item-subtitle>
                </v-list-item>
              </v-list>
              <v-card v-else variant="outlined" class="text-center pa-4">
                <p class="text-grey-lighten-1">Nenhuma subscription encontrada</p>
              </v-card>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="topicDetailsDialog = false">Fechar</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">Confirmar Exclusão</v-card-title>
        <v-card-text>
          Tem certeza que deseja excluir o tópico <strong>{{ getTopicName(topicToDelete) }}</strong>?
          Esta ação não pode ser desfeita.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="deleteDialog = false">Cancelar</v-btn>
          <v-btn color="error" @click="deleteTopic" :loading="deleting">Excluir</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAppStore } from '@/stores/app'
import { storeToRefs } from 'pinia'
import TitleNameWithTooltip from '@/components/TitleNameWithTooltip.vue'
import { 
  ListTopicsCommand, 
  CreateTopicCommand, 
  DeleteTopicCommand, 
  PublishCommand, 
  ListSubscriptionsByTopicCommand 
} from '@aws-sdk/client-sns'

const appStore = useAppStore()
const { sns } = storeToRefs(appStore)

// Reactive data
const topics = ref([])
const searchQuery = ref('')
const loading = ref(false)
const creating = ref(false)
const publishing = ref(false)
const deleting = ref(false)

// Dialog states
const createTopicDialog = ref(false)
const publishDialog = ref(false)
const topicDetailsDialog = ref(false)
const deleteDialog = ref(false)

// Form data
const newTopicName = ref('')
const messageSubject = ref('')
const messageContent = ref('')
const topicToDelete = ref('')
const topicToPublish = ref('')
const selectedTopic = ref(null)
const subscriptions = ref([])

// Validation rules
const rules = {
  required: value => !!value || 'Campo obrigatório'
}

// Computed properties
const filteredTopics = computed(() => {
  if (!searchQuery.value) return topics.value
  return topics.value.filter(topic => 
    getTopicName(topic.TopicArn).toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// Helper functions
const getTopicName = (topicArn) => {
  return topicArn ? topicArn.split(':').pop() : ''
}

// SNS API functions
const loadTopics = async () => {
  try {
    loading.value = true
    const response = await sns.value.send(new ListTopicsCommand({}))
    
    topics.value = response.Topics || []
  } catch (error) {
    console.error('Erro ao carregar tópicos:', error)
    appStore.showSnackbar('Erro ao carregar tópicos SNS', 'error')
  } finally {
    loading.value = false
  }
}

const createTopic = async () => {
  if (!newTopicName.value) return
  
  try {
    creating.value = true
    const response = await sns.value.send(new CreateTopicCommand({
      Name: newTopicName.value
    }))
    
    appStore.showSnackbar('Tópico criado com sucesso!', 'success')
    createTopicDialog.value = false
    newTopicName.value = ''
    await loadTopics()
  } catch (error) {
    console.error('Erro ao criar tópico:', error)
    appStore.showSnackbar('Erro ao criar tópico', 'error')
  } finally {
    creating.value = false
  }
}

const openTopic = async (topicArn) => {
  try {
    selectedTopic.value = { TopicArn: topicArn }
    
    // Load subscriptions
    const response = await sns.value.send(new ListSubscriptionsByTopicCommand({
      TopicArn: topicArn
    }))
    
    subscriptions.value = response.Subscriptions || []
    topicDetailsDialog.value = true
  } catch (error) {
    console.error('Erro ao carregar detalhes do tópico:', error)
    appStore.showSnackbar('Erro ao carregar detalhes do tópico', 'error')
  }
}

const publishToTopic = (topicArn) => {
  topicToPublish.value = topicArn
  messageSubject.value = ''
  messageContent.value = ''
  publishDialog.value = true
}

const publishMessage = async () => {
  if (!messageContent.value) return
  
  try {
    publishing.value = true
    const payload = {
      TopicArn: topicToPublish.value,
      Message: messageContent.value
    }
    
    if (messageSubject.value) {
      payload.Subject = messageSubject.value
    }
    
    await sns.value.send(new PublishCommand(payload))
    
    appStore.showSnackbar('Mensagem publicada com sucesso!', 'success')
    publishDialog.value = false
  } catch (error) {
    console.error('Erro ao publicar mensagem:', error)
    appStore.showSnackbar('Erro ao publicar mensagem', 'error')
  } finally {
    publishing.value = false
  }
}

const confirmDeleteTopic = (topicArn) => {
  topicToDelete.value = topicArn
  deleteDialog.value = true
}

const deleteTopic = async () => {
  try {
    deleting.value = true
    await sns.value.send(new DeleteTopicCommand({
      TopicArn: topicToDelete.value
    }))
    
    appStore.showSnackbar('Tópico excluído com sucesso!', 'success')
    deleteDialog.value = false
    await loadTopics()
  } catch (error) {
    console.error('Erro ao excluir tópico:', error)
    appStore.showSnackbar('Erro ao excluir tópico', 'error')
  } finally {
    deleting.value = false
  }
}

// Lifecycle
onMounted(() => {
  loadTopics()
})
</script>

<style scoped>
.text-break {
  word-break: break-all;
}
</style>
