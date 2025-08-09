<template>
  <div class="dashboard">
    <!-- Header -->
    <v-row class="mb-4">
      <v-col cols="12">
        <div class="d-flex align-center justify-space-between">
          <div class="d-flex align-center">
            <v-icon class="mr-2" size="24">mdi-view-dashboard-outline</v-icon>
            <h1 class="text-h5 mb-0">Dashboard</h1>
          </div>
          <v-btn size="small" prepend-icon="mdi-refresh" variant="text" @click="refreshAll">
            Atualizar
          </v-btn>
        </div>
      </v-col>
    </v-row>

    <!-- KPI / Services Grid -->
    <v-row>
      <v-col cols="12" md="6" lg="4" xl="3" v-for="service in services" :key="service.name">
        <v-card class="kpi-card" elevation="1" rounded="lg" variant="elevated" @click="service.status === 'active' ? navigateToService(service.route): null">
          <v-card-text class="py-4">
            <div class="d-flex align-start justify-space-between mb-4">
              <div class="d-flex align-center">
                <div class="icon-chip mr-3">
                  <v-icon size="22">{{ service.icon }}</v-icon>
                </div>
                <div>
                  <div class="text-subtitle-1 font-weight-medium">{{ service.name }}</div>
                  <v-chip :color="service.status === 'active' ? 'success' : 'warning'" size="x-small" variant="flat" class="mt-1">
                    {{ service.status === 'active' ? 'Ativo' : 'Inativo' }}
                  </v-chip>
                </div>
              </div>
              <v-icon v-if="service.status === 'active'" size="18" color="success">mdi-check-circle</v-icon>
              <v-icon v-else size="18" color="warning">mdi-alert-circle</v-icon>
            </div>

            <div class="stats">
              <template v-if="service.status === 'active'">
                <template v-if="service.stats">
                  <div v-for="(value, key) in service.stats" :key="key" class="stat-line">
                    <span class="stat-key">{{ key }}</span>
                    <span class="stat-value">{{ value }}</span>
                  </div>
                </template>
                <template v-else>
                  <v-skeleton-loader type="text" class="mb-1"></v-skeleton-loader>
                  <v-skeleton-loader type="text" width="40%"></v-skeleton-loader>
                </template>
              </template>
              <template v-else>
                <div class="text-caption text-medium-emphasis">Serviço não disponível</div>
              </template>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Quick Actions -->
    <v-row class="mt-6">
      <v-col cols="12">
        <v-card elevation="1" rounded="lg">
          <v-card-title class="py-3 d-flex align-center">
            <v-icon class="mr-2">mdi-flash</v-icon>
            Ações rápidas
          </v-card-title>
          <v-divider></v-divider>
          <v-card-text class="pt-5">
            <v-row>
              <v-col cols="12" sm="6" md="4" lg="3" v-for="action in quickActions" :key="action.title">
                <v-btn
                  :color="action.color"
                  variant="tonal"
                  block
                  class="action-btn mb-3"
                  size="large"
                  @click="executeQuickAction(action)"
                  :loading="action.loading"
                >
                  <v-icon class="mr-2">{{ action.icon }}</v-icon>
                  {{ action.title }}
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Confirmation Dialog -->
    <v-dialog v-model="confirmDialog.show" max-width="420">
      <v-card rounded="lg">
        <v-card-title class="text-h6">
          {{ confirmDialog.title }}
        </v-card-title>
        <v-card-text>
          {{ confirmDialog.text }}
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="confirmDialog.show = false">Cancelar</v-btn>
          <v-btn :color="confirmDialog.color" @click="confirmDialog.confirm" :loading="confirmDialog.loading">Confirmar</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { storeToRefs } from 'pinia'

// AWS SDK v3 Commands
import { ListTablesCommand, DeleteTableCommand } from '@aws-sdk/client-dynamodb'
import { ListStreamsCommand } from '@aws-sdk/client-kinesis'
import { ListKeysCommand } from '@aws-sdk/client-kms'
import { ListFunctionsCommand } from '@aws-sdk/client-lambda'
import { ListBucketsCommand, ListObjectsV2Command, DeleteObjectsCommand, DeleteBucketCommand } from '@aws-sdk/client-s3'
import { ListTopicsCommand } from '@aws-sdk/client-sns'
import { ListQueuesCommand, DeleteQueueCommand } from '@aws-sdk/client-sqs'
import { ListIdentitiesCommand } from '@aws-sdk/client-ses'
import { deleteSesMessage, getSesData } from '@/utils/api.js'

const router = useRouter()
const appStore = useAppStore()
const { s3, ses, sns, sqs, dynamodb, lambda, kinesis, kms } = storeToRefs(appStore)

const services = ref([
  { name: 'DynamoDB', icon: 'mdi-table', route: '/dynamodb', status: 'inactive', stats: null },
  { name: 'Kinesis', icon: 'mdi-view-stream', route: '/kinesis', status: 'inactive', stats: null },
  { name: 'KMS', icon: 'mdi-key', route: '/kms', status: 'inactive', stats: null },
  { name: 'Lambda', icon: 'mdi-function', route: '/lambda', status: 'inactive', stats: null },
  { name: 'S3', icon: 'mdi-folder-multiple', route: '/s3', status: 'inactive', stats: null },
  { name: 'SES', icon: 'mdi-email', route: '/ses', status: 'inactive', stats: null },
  { name: 'SNS', icon: 'mdi-forum', route: '/sns', status: 'inactive', stats: null },
  { name: 'SQS', icon: 'mdi-format-list-bulleted', route: '/sqs', status: 'inactive', stats: null },
])

const quickActions = ref([
  { title: 'Limpar S3', icon: 'mdi-delete-sweep', color: 'error', action: 'clearS3', loading: false },
  { title: 'Limpar SQS', icon: 'mdi-broom', color: 'warning', action: 'clearSQS', loading: false },
  { title: 'Limpar DynamoDB', icon: 'mdi-table-remove', color: 'orange', action: 'clearDynamoDB', loading: false },
  { title: 'Limpar SES', icon: 'mdi-email-remove', color: 'error', action: 'clearSES', loading: false },
  { title: 'Atualizar Status', icon: 'mdi-refresh', color: 'primary', action: 'refreshStatus', loading: false },
])

const confirmDialog = ref({ show: false, title: '', text: '', color: 'primary', confirm: null, loading: false })

const navigateToService = (route) => router.push(route)
const refreshAll = () => {
  const refresh = quickActions.value.find(a => a.action === 'refreshStatus')
  if (refresh) executeQuickAction(refresh)
}

const loadServiceStats = async () => {
  try {
    await Promise.all([
      loadDynamoStats(),
      loadKinesisStats(),
      loadKmsStats(),
      loadLambdaStats(),
      loadS3Stats(),
      loadSesStats(),
      loadSnsStats(),
      loadSqsStats(),
    ])
  } catch (error) {
    console.error('Error loading service stats:', error)
    appStore.showSnackbar('Erro ao carregar estatísticas dos serviços', 'error')
  }
}

const loadDynamoStats = async () => {
  try {
    const res = await dynamodb.value.send(new ListTablesCommand({}))
    services.value[0].status = 'active'
    services.value[0].stats = { 'Tabelas': res.TableNames ? res.TableNames.length : 0 }
  } catch (e) {
    console.error('DynamoDB error:', e)
    services.value[0].status = 'inactive'
  }
}

const loadKinesisStats = async () => {
  try {
    const res = await kinesis.value.send(new ListStreamsCommand({}))
    services.value[1].status = 'active'
    services.value[1].stats = { 'Streams': res.StreamNames ? res.StreamNames.length : 0 }
  } catch (e) {
    console.error('Kinesis error:', e)
    services.value[1].status = 'inactive'
  }
}

const loadKmsStats = async () => {
  try {
    const res = await kms.value.send(new ListKeysCommand({}))
    services.value[2].status = 'active'
    services.value[2].stats = { 'Chaves': res.Keys ? res.Keys.length : 0 }
  } catch (e) {
    console.error('KMS error:', e)
    services.value[2].status = 'inactive'
  }
}

const loadLambdaStats = async () => {
  try {
    const res = await lambda.value.send(new ListFunctionsCommand({}))
    services.value[3].status = 'active'
    services.value[3].stats = { 'Funções': res.Functions ? res.Functions.length : 0 }
  } catch (e) {
    console.error('Lambda error:', e)
    services.value[3].status = 'inactive'
  }
}

const loadS3Stats = async () => {
  try {
    const res = await s3.value.send(new ListBucketsCommand({}))
    services.value[4].status = 'active'
    services.value[4].stats = { 'Buckets': res.Buckets ? res.Buckets.length : 0 }
  } catch (e) {
    console.error('S3 error:', e)
    services.value[4].status = 'inactive'
  }
}

const loadSesStats = async () => {
  try {
    const identities = await ses.value.send(new ListIdentitiesCommand({ IdentityType: 'EmailAddress' }))
    const data = await getSesData()
    services.value[5].status = 'active'
    services.value[5].stats = { 'E-mails': data.messages ? data.messages.length : 0, 'Identidades': identities.Identities ? identities.Identities.length : 0 }
  } catch (e) {
    console.error('SES error:', e)
    services.value[5].status = 'inactive'
  }
}

const loadSnsStats = async () => {
  try {
    const res = await sns.value.send(new ListTopicsCommand({}))
    services.value[6].status = 'active'
    services.value[6].stats = { 'Tópicos': res.Topics ? res.Topics.length : 0 }
  } catch (e) {
    console.error('SNS error:', e)
    services.value[6].status = 'inactive'
  }
}

const loadSqsStats = async () => {
  try {
    const res = await sqs.value.send(new ListQueuesCommand({}))
    services.value[7].status = 'active'
    services.value[7].stats = { 'Filas': res.QueueUrls ? res.QueueUrls.length : 0 }
  } catch (e) {
    console.error('SQS error:', e)
    services.value[7].status = 'inactive'
  }
}

const executeQuickAction = (action) => {
  if (action.action === 'refreshStatus') {
    action.loading = true
    loadServiceStats().finally(() => {
      action.loading = false
      appStore.showSnackbar('Status atualizado!', 'success')
    })
    return
  }

  let title, text, color
  switch (action.action) {
    case 'clearS3':
      title = 'Limpar todos os buckets S3'
      text = 'Esta ação irá deletar todos os buckets e objetos S3. Deseja continuar?'
      color = 'error'
      break
    case 'clearSQS':
      title = 'Limpar todas as filas SQS'
      text = 'Esta ação irá deletar todas as filas SQS. Deseja continuar?'
      color = 'warning'
      break
    case 'clearSES':
      title = 'Limpar todos os emails SES'
      text = 'Esta ação irá deletar todos os emails do SES. Deseja continuar?'
      color = 'error'
      break
    case 'clearDynamoDB':
      title = 'Limpar todas as tabelas DynamoDB'
      text = 'Esta ação irá deletar todas as tabelas DynamoDB. Deseja continuar?'
      color = 'orange'
      break
  }

  confirmDialog.value = { show: true, title, text, color, confirm: () => performClearAction(action), loading: false }
}

const performClearAction = async (action) => {
  confirmDialog.value.loading = true
  action.loading = true
  try {
    switch (action.action) {
      case 'clearS3':
        await clearAllS3Buckets()
        break
      case 'clearSQS':
        await clearAllSQSQueues()
        break
      case 'clearSES':
        await clearAllSESData()
        break
      case 'clearDynamoDB':
        await clearAllDynamoDBTables()
        break
    }
    appStore.showSnackbar(`${action.title} executado com sucesso!`, 'success')
    await loadServiceStats()
  } catch (error) {
    console.error(`Error in ${action.action}:`, error)
    appStore.showSnackbar(`Erro ao executar ${action.title}`, 'error')
  } finally {
    confirmDialog.value.show = false
    confirmDialog.value.loading = false
    action.loading = false
  }
}

const clearAllS3Buckets = async () => {
  const buckets = await s3.value.send(new ListBucketsCommand({}))
  for (const bucket of buckets.Buckets) {
    const objects = await s3.value.send(new ListObjectsV2Command({ Bucket: bucket.Name }))
    if (objects.Contents && objects.Contents.length > 0) {
      const deleteParams = { Bucket: bucket.Name, Delete: { Objects: objects.Contents.map(obj => ({ Key: obj.Key })) } }
      await s3.value.send(new DeleteObjectsCommand(deleteParams))
    }
    await s3.value.send(new DeleteBucketCommand({ Bucket: bucket.Name }))
  }
}

const clearAllSQSQueues = async () => {
  const queues = await sqs.value.send(new ListQueuesCommand({}))
  if (queues.QueueUrls) {
    for (const queueUrl of queues.QueueUrls) {
      await sqs.value.send(new DeleteQueueCommand({ QueueUrl: queueUrl }))
    }
  }
}

const clearAllSESData = async () => {
  await deleteSesMessage()
}

const clearAllDynamoDBTables = async () => {
  const tables = await dynamodb.value.send(new ListTablesCommand({}))
  for (const tableName of tables.TableNames) {
    await dynamodb.value.send(new DeleteTableCommand({ TableName: tableName }))
  }
}

onMounted(() => {
  loadServiceStats()
})
</script>

<style scoped>
.kpi-card {
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}
.kpi-card:hover {
  transform: translateY(-2px);
}

.icon-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.06);
}

:deep(.v-theme--dark) .icon-chip {
  background: rgba(255, 255, 255, 0.08);
}

.stat-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 0;
  border-bottom: 1px dashed rgba(0, 0, 0, 0.06);
}
.stat-line:last-child {
  border-bottom: none;
}
.stat-key {
  color: rgba(0, 0, 0, 0.6);
}
:deep(.v-theme--dark) .stat-key {
  color: rgba(255, 255, 255, 0.6);
}
.stat-value {
  font-weight: 600;
}

.action-btn {
  height: 44px;
}
</style>
