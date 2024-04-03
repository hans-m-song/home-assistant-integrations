{{/*
Common labels
*/}}
{{- define "home-assistant-integrations.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{ include "home-assistant-integrations.selectors" . }}
{{- end }}

{{/*
Common selectors
*/}}
{{- define "home-assistant-integrations.selectors" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
