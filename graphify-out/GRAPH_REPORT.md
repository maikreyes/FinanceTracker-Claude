# Graph Report - FinanceTracker  (2026-09-27)

## Corpus Check
- 157 files · ~55,668 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 2 file(s) not represented in the graph (top: .example 1, (none) 1)

## Summary
- 743 nodes · 1895 edges · 39 communities (28 shown, 11 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 156 edges (avg confidence: 0.85)
- Token cost: 484,910 input · 0 output

## Community Hubs (Navigation)
- Webhook Entry & Error Types
- Test Suite Coverage
- Bitacora de Features
- Notion & Reply Handlers
- Documentacion de Arquitectura
- Bootstrap & Config
- Telegram Messenger Port
- Test Fakes & Harness
- Legacy Bot Handlers
- Upstash Store
- Webhook Dispatch & ChatState
- Telegram Command Handlers
- ChatState Expiry & Clock
- Dispatch & State Persistence
- Guided Wizard Flow
- Bot Entry & Notion Schema
- Monthly Summary Feature
- Balance & Summary Rationale
- Confirmation & Error Handling
- Access Whitelist & Summary Specs
- Telegram Services Tests
- Receipt Photo & Type Specs
- Groq Receipt Interpretation
- Category Resolution
- Receipt Upload API
- Clean Architecture Setup (FDC)
- Core Transaction Domain (Legacy)
- Notion File Upload Services
- Telegram Bot Commands Setup
- Fake Category Resolver
- Vercel Deployment Config
- Receipt Interpretation Error Type
- Write Failure Error Type
- ValidAmount Fix
- WithoutURL Token Guard
- Go Module Root
- Feature Template Plan
- Feature Template Spec
- Feature Template Tasks

## God Nodes (most connected - your core abstractions)
1. `newHarness()` - 42 edges
2. `ChatState` - 34 edges
3. `NewServices()` - 26 edges
4. `Message` - 16 edges
5. `send()` - 14 edges
6. `Transaction` - 14 edges
7. `harness` - 13 edges
8. `tipoCallbackData()` - 13 edges
9. `Result` - 13 edges
10. `usecase.Registrar (Register/RegisterWithPhoto/RegisterFromPhotoAI/RegisterFromFields)` - 13 edges

## Surprising Connections (you probably didn't know these)
- `Webhook siempre responde 200 aunque falle el despacho interno` --rationale_for--> `Handler()`  [EXTRACTED]
  specs/features/012-migracion-arquitectura-vercel/plan.md → api/webhook.go
- `Old Tech Stack Doc (domain/usecase/infrastructure, not updated to pkg/ pattern)` --conceptually_related_to--> `AGENT.md Project Structure (api/, cmd/, pkg/<feature>/{handler,model,services}, pkg/ports)`  [AMBIGUOUS]
  specs/constitution/tech-stack.md → AGENT.md
- `Handler()` --calls--> `bot.Dispatch`  [EXTRACTED]
  api/webhook.go → specs/features/012-migracion-arquitectura-vercel/plan.md
- `Handler()` --shares_data_with--> `upstash.Store (Redis-backed Store)`  [EXTRACTED]
  api/webhook.go → specs/features/012-migracion-arquitectura-vercel/plan.md
- `Plan 012: Migración de arquitectura para Vercel` --references--> `Handler()`  [EXTRACTED]
  specs/features/012-migracion-arquitectura-vercel/plan.md → api/webhook.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Patrón puertos/adaptadores: usecase.ports desacoplado de implementaciones concretas de Notion/Telegram/Groq** — internal_usecase_ports, internal_infrastructure_notion_client, internal_infrastructure_notion_category_resolver, internal_infrastructure_notion_summary, internal_infrastructure_notion_receipt_uploader, internal_infrastructure_notion_latest, internal_infrastructure_groq_client, internal_infrastructure_telegram_client [EXTRACTED 1.00]
- **Roadmap secuencial de features 001-012 del Finance Tracker** — concept_feature_001_registrar_movimiento, concept_feature_002_transporte_telegram, concept_feature_003_restringir_acceso_bot, concept_feature_004_confirmacion_manejo_errores, concept_feature_005_resumen_gastos, concept_feature_006_foto_recibo, concept_feature_007_interpretar_recibo_ia, concept_feature_008_botones_tipo, concept_feature_009_balance_global_resumen, concept_feature_010_comandos_flujo_guiado, concept_feature_011_comandos_ayuda_ultimo_categorias, concept_feature_012_migracion_vercel_webhook [EXTRACTED 1.00]
- **Mecanismo de estado pendiente por chat_id que motivó el refactor a internal/bot** — concept_pending_state_in_memory_pattern, internal_bot_state, cmd_bot_main, concept_feature_012_migracion_vercel_webhook [INFERRED 0.85]
- **Feature+Ports Architecture Convention** — claude_rules_codigo_architecture_pattern, agent_project_structure, claude_history_2026_09_26_arquitectura_patron_triggo_y_vercel_triggo_pattern [INFERRED 0.85]
- **Receipt Photo Capture + AI Interpretation Pipeline** — claude_history_2026_09_20_spec006_adjuntar_recibo_foto_file_upload_api, claude_history_2026_09_20_spec006_adjuntar_recibo_foto_registerwithphoto, claude_history_2026_09_20_spec007_interpretar_recibo_ia_feature007, specs_constitution_roadmap_feature008 [INFERRED 0.85]
- **Thermos Review Corrections to Feature 012 Step 1** — claude_history_2026_09_26_bot_correcciones_revision_thermos_chatstate, claude_history_2026_09_26_bot_correcciones_revision_thermos_withouturl, claude_history_2026_09_26_bot_correcciones_revision_thermos_validamount, claude_history_2026_09_26_bot_correcciones_revision_thermos_dispatch_split [EXTRACTED 1.00]
- **Category resolution ports pipeline (parse -> interface -> Notion impl -> orchestrator)** — internal_usecase_parse_transaction_parsemessage, internal_usecase_ports_categoryresolver, internal_infrastructure_notion_category_resolver_categoryresolver, internal_usecase_register_transaction_registrar [INFERRED 0.85]
- **Receipt photo upload pipeline (Telegram download -> ReceiptUploader -> Notion File Upload API -> Registrar)** — internal_infrastructure_telegram_client_bot, internal_usecase_ports_receiptuploader, internal_infrastructure_notion_receipt_uploader_receiptuploader, internal_usecase_register_transaction_registrar, notion_file_upload_api [INFERRED 0.85]
- **AI receipt interpretation with button-based type confirmation** — internal_usecase_ports_receiptinterpreter, internal_infrastructure_groq_client_client, internal_usecase_register_transaction_registrar, internal_infrastructure_telegram_client_callbackquery [INFERRED 0.85]
- **Template spec/plan/tasks instanciado por las features 009-012** — specs_features_nnn_nombre_feature_spec, specs_features_nnn_nombre_feature_plan, specs_features_nnn_nombre_feature_tasks, specs_features_009_balance_global_resumen_spec, specs_features_010_comandos_flujo_guiado_spec, specs_features_011_comandos_ayuda_ultimo_categorias_spec, specs_features_012_migracion_arquitectura_vercel_spec [INFERRED 0.85]
- **Evolución del estado de conversación: mapas en memoria -> ChatState + Store (memoria/Upstash)** — cmd_bot_main_pendingconversations, internal_bot_state_chatstate, internal_bot_state_store, internal_bot_memorystore_memorystore, internal_infrastructure_upstash_store [INFERRED 0.85]
- **Pipeline de despliegue webhook en Vercel** — api_webhook_handler, internal_bot_dispatch_dispatch, internal_infrastructure_upstash_store, internal_infrastructure_telegram_client_setwebhook, cmd_setwebhook_main [INFERRED 0.85]

## Communities (39 total, 11 thin omitted)

### Community 0 - "Webhook Entry & Error Types"
Cohesion: 0.08
Nodes (41): ErrInvalidAmount, ErrInvalidFormat, ErrReceiptAmountNotFound, ErrUnrecognizedCategory, ErrUnrecognizedPaymentMethod, ErrUnrecognizedType, go_pkg_bytes, go_pkg_context (+33 more)

### Community 1 - "Test Suite Coverage"
Cohesion: 0.05
Nodes (88): net/http/httptest.ResponseRecorder, testing.T, TestCommands_IncluyenResumen(), TestNew_CablaTodosLosPuertos(), TestCheckWebhook(), TestCategoryResolver_Create_CategoriaExistenteNoDuplica(), TestCategoryResolver_Create_CategoriaNueva(), TestClient_FindLatest_ConCategoria() (+80 more)

### Community 2 - "Bitacora de Features"
Cohesion: 0.07
Nodes (56): Doc: Notion MCP schema mapping, Doc: Notion real request validation, Doc: Feature 001 implementada, Doc: Feature 002 implementada, Doc: Feature 003 implementada, Doc: Feature 004 implementada, Doc: Feature 005 implementada, Doc: Feature 006 implementada (+48 more)

### Community 3 - "Notion & Reply Handlers"
Cohesion: 0.05
Nodes (35): fakeFinder, fakeWriter, Services, mapTransactionToProperties(), Services, TransactionServices, ultimoReply(), TestUltimoReply() (+27 more)

### Community 4 - "Documentacion de Arquitectura"
Cohesion: 0.05
Nodes (52): Data Flow: Telegram -> handler -> transaction services -> notion/groq services, AGENT.md Project Structure (api/, cmd/, pkg/<feature>/{handler,model,services}, pkg/ports), CLAUDE.md Bitácora Index (resumen + índice + última entrada completa), mapTransactionToProperties (notion client Type mapping), domain.Transaction.Type (TransactionType), Type Property Added to Expenses (Ingreso/Egreso), CategoryResolver interface (live Notion resolution), Feature 001: Registrar movimiento por mensaje de Telegram (+44 more)

### Community 5 - "Bootstrap & Config"
Cohesion: 0.06
Nodes (34): initApp(), main(), runPollingLoop(), main(), TestValidateWebhookURL(), validateWebhookURL(), net/http.Client, net/url.Values (+26 more)

### Community 6 - "Telegram Messenger Port"
Cohesion: 0.10
Nodes (11): context.Context, fakeMessenger, Messenger, Handler, editMsg(), InlineButton, Services, Services (+3 more)

### Community 7 - "Test Fakes & Harness"
Cohesion: 0.07
Nodes (13): fakeCategories, fakeInterpreter, harness, jsonStore, buildPrompt(), Services, ReceiptInterpreter, ExpenseWriter (+5 more)

### Community 8 - "Legacy Bot Handlers"
Cohesion: 0.08
Nodes (29): handleCategoriaCallback, handleCategoriasCommand, handleCategoryAddReply, handleConversationReply, handleStartConversation, handleUltimo, matchesCommand, parseStartCommand (+21 more)

### Community 9 - "Upstash Store"
Cohesion: 0.16
Nodes (18): encoding/json.RawMessage, chatKey(), decodeState(), newUpstash(), newFakeRedis(), TestUpstash_Delete(), TestUpstash_ErrorDeLaAPI(), TestUpstash_GetSinEstadoDevuelveCero() (+10 more)

### Community 10 - "Webhook Dispatch & ChatState"
Cohesion: 0.12
Nodes (21): Handler(), TestHandler_ConfiguracionInvalidaResponde500SinFiltrarDetalle(), runPollingLoop, net/http.Request, net/http.ResponseWriter, bot.Deps struct, bot.Dispatch, ChatState struct (+13 more)

### Community 11 - "Telegram Command Handlers"
Cohesion: 0.18
Nodes (12): Handler, send(), Handler, errorReply(), CallbackQuery, Message, Update, Services (+4 more)

### Community 12 - "ChatState Expiry & Clock"
Cohesion: 0.23
Nodes (7): PendingKind, time.Time, clone(), ChatState, NewPhotoPending(), nowBogota(), Memory

### Community 13 - "Dispatch & State Persistence"
Cohesion: 0.21
Nodes (7): matchesCommand(), normalizeCommand(), TestIsAllowed(), TestMatchesCommand(), TestNormalizeCommand(), Handler, isAllowed()

### Community 14 - "Guided Wizard Flow"
Cohesion: 0.23
Nodes (10): WizardStep, parseStartCommand(), successReply(), TestSuccessReply_SinCategoria(), advanceWizard(), Handler, skipDash(), Wizard (+2 more)

### Community 15 - "Bot Entry & Notion Schema"
Cohesion: 0.20
Nodes (10): cmd/bot/main.go entrypoint (wiring + polling loop), telegram.Bot client, Notion data source: Category, Notion data source: Expenses, Rationale: net/http long-polling, no Telegram SDK, no webhook, Rationale: typed errors in usecase, chat copy owned by transport layer, 001 Registrar movimiento — Spec, 002 Transporte Telegram — Plan (+2 more)

### Community 16 - "Monthly Summary Feature"
Cohesion: 0.24
Nodes (6): fakeSummarizer, Services, resumenReply(), TestResumenReply(), TestResumenReply_SinMovimientos(), Monthly

### Community 17 - "Balance & Summary Rationale"
Cohesion: 0.28
Nodes (9): resumenReply, notion.Client.SumThisMonth, MonthlySummary struct + Balance(), Backfill por curl directo, no comando de la app, Balance = Ingreso - Egreso (corregido tras aclaración del usuario), Rationale: client-side sum instead of relying on Add to Month formula, Plan 009: Balance global + backfill, Spec 009: Balance global + backfill de Tipo (+1 more)

### Community 18 - "Confirmation & Error Handling"
Cohesion: 0.22
Nodes (9): ErrUnrecognizedCategory, ErrUnrecognizedPaymentMethod, ErrUnrecognizedType, ErrWriteFailed, RegisterFromPhotoAI: unmatched category/payment method does not block registration, Rationale: AI-suggested category/payment-method that don't match are left empty, never block registration, 004 Confirmación y manejo de errores — Plan, 004 Confirmación y manejo de errores — Spec (+1 more)

### Community 19 - "Access Whitelist & Summary Specs"
Cohesion: 0.25
Nodes (8): isAllowed + TELEGRAM_ALLOWED_CHAT_IDS whitelist check, Rationale: fail-closed whitelist (no chats allowed without config), 003 Restringir acceso al bot — Plan, 003 Restringir acceso al bot — Spec, 003 Restringir acceso al bot — Tasks, 005 Consultar resumen de gastos — Plan, 005 Consultar resumen de gastos — Spec, 005 Consultar resumen de gastos — Tasks

### Community 20 - "Telegram Services Tests"
Cohesion: 0.46
Nodes (7): newServices(), closedServerURL(), TestCall_ErrorDeRedNoFiltraElToken(), TestDeleteWebhook_DescartaPendientesSoloSiSePide(), TestDownloadFile_ErrorDeRedNoFiltraElToken(), TestDownloadPhoto_ResuelveFileIDYDescarga(), TestSetWebhook_EnviaUrlYSecreto()

### Community 21 - "Receipt Photo & Type Specs"
Cohesion: 0.33
Nodes (7): pendingPhotos map[int64][]byte (in-memory photo staging), CallbackQuery / inline keyboard mechanism, 006 Adjuntar foto de recibo — Spec, 007 Interpretar recibo con IA — Spec, 008 Preguntar tipo con botones — Plan, 008 Preguntar tipo con botones — Spec, 008 Preguntar tipo con botones — Tasks

### Community 22 - "Groq Receipt Interpretation"
Cohesion: 0.33
Nodes (7): groq.Client (vision interpretation), ErrReceiptAmountNotFound, ErrReceiptInterpretationFailed, ReceiptInterpreter interface, Rationale: Groq chosen over Gemini (billing link requirement), 007 Interpretar recibo con IA — Plan, 007 Interpretar recibo con IA — Tasks

### Community 23 - "Category Resolution"
Cohesion: 0.40
Nodes (6): notion.CategoryResolver implementation, ParseMessage (positional parser), CategoryResolver interface, Rationale: categories resolved live against Notion, never hardcoded, 001 Registrar movimiento — Plan, 001 Registrar movimiento — Tasks

### Community 24 - "Receipt Upload API"
Cohesion: 0.40
Nodes (6): notion.receipt_uploader implementation, ReceiptUploader interface, Notion File Upload API (two-step upload), Rationale: photo uploaded natively to Notion, not linked (Telegram links expire), 006 Adjuntar foto de recibo — Plan, 006 Adjuntar foto de recibo — Tasks

### Community 25 - "Clean Architecture Setup (FDC)"
Cohesion: 0.40
Nodes (5): Doc: Documentación adaptar patrón FDC, Doc: Setup andamiaje Go clean architecture, Doc: specs/ replicar carpeta FDC, Clean Architecture (domain/usecase/infrastructure), moveeng/FDC/DESARROLLO (proyecto de referencia para memoria/specs)

### Community 26 - "Core Transaction Domain (Legacy)"
Cohesion: 0.40
Nodes (5): domain.Transaction, notion.Client, MonthlySummarizer interface, Registrar orchestrator (Register/RegisterWithPhoto/RegisterFromPhotoAI), Rationale: manual caption takes priority over AI interpretation

### Community 30 - "Vercel Deployment Config"
Cohesion: 0.50
Nodes (3): maxDuration, functions, api/webhook.go

## Ambiguous Edges - Review These
- `AGENT.md Project Structure (api/, cmd/, pkg/<feature>/{handler,model,services}, pkg/ports)` → `Old Tech Stack Doc (domain/usecase/infrastructure, not updated to pkg/ pattern)`  [AMBIGUOUS]
  specs/constitution/tech-stack.md · relation: conceptually_related_to

## Knowledge Gaps
- **67 isolated node(s):** `finance-tracker`, `Services`, `chatCompletionResponse`, `receiptData`, `categoryQueryResponse` (+62 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 145 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **11 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `AGENT.md Project Structure (api/, cmd/, pkg/<feature>/{handler,model,services}, pkg/ports)` and `Old Tech Stack Doc (domain/usecase/infrastructure, not updated to pkg/ pattern)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **Why does `Plan 012: Migración de arquitectura para Vercel` connect `Webhook Dispatch & ChatState` to `Webhook Entry & Error Types`, `Legacy Bot Handlers`, `Bitacora de Features`?**
  _High betweenness centrality (0.181) - this node is a cross-community bridge._
- **Why does `internal/bot/memorystore.Store (in-memory)` connect `Bitacora de Features` to `Webhook Dispatch & ChatState`?**
  _High betweenness centrality (0.110) - this node is a cross-community bridge._
- **Are the 37 inferred relationships involving `newHarness()` (e.g. with `TestDispatch_AgregarCategoriaCancelaFlujoGuiado()` and `TestDispatch_BotonesDeFotoVencidosYaNoAplican()`) actually correct?**
  _`newHarness()` has 37 INFERRED edges - model-reasoned connections that need verification._
- **Are the 18 inferred relationships involving `NewServices()` (e.g. with `TestRegistrar_RegisterFromFields_MontoInvalido()` and `TestRegistrar_RegisterFromPhotoAI_MontoInvalidoEsMontoNoEncontrado()`) actually correct?**
  _`NewServices()` has 18 INFERRED edges - model-reasoned connections that need verification._
- **What connects `finance-tracker`, `Services`, `chatCompletionResponse` to the rest of the system?**
  _67 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Webhook Entry & Error Types` be split into smaller, more focused modules?**
  _Cohesion score 0.08104331625523988 - nodes in this community are weakly interconnected._