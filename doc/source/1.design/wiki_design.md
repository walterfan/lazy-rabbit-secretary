# Wiki Module Design Document

## Overview

The Wiki module provides a collaborative knowledge base system with full Markdown support, version control, and user management. It allows users to create, edit, and manage wiki pages with features like categories, tags, page protection, and revision history.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Backend Design](#backend-design)
3. [Frontend Design](#frontend-design)
4. [API Design](#api-design)
5. [Database Schema](#database-schema)
6. [Security Model](#security-model)
7. [Markdown Rendering](#markdown-rendering)
8. [Version Control](#version-control)
9. [User Experience](#user-experience)
10. [Internationalization](#internationalization)

## Architecture Overview

The wiki module follows a layered architecture pattern:

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend Layer                           │
├─────────────────────────────────────────────────────────────┤
│  Views: WikiView, WikiPageView, WikiEditView               │
│  Components: WikiSearchBar, WikiSidebar, WikiPageList      │
│  Store: WikiStore (Pinia)                                  │
│  Router: /wiki, /wiki/page/:slug, /wiki/edit/:slug?       │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Layer                                │
├─────────────────────────────────────────────────────────────┤
│  Public Routes: /api/v1/wiki/*                              │
│  Admin Routes: /api/v1/admin/wiki/*                        │
│  Middleware: OptionalAuth, RequiredAuth                    │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                  Service Layer                              │
├─────────────────────────────────────────────────────────────┤
│  WikiService: Business logic, validation, permissions      │
│  WikiRepository: Database operations                       │
│  Models: WikiPage, WikiRevision                            │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                  Database Layer                             │
├─────────────────────────────────────────────────────────────┤
│  Tables: wiki_pages, wiki_revisions                        │
│  Indexes: slug, status, type, categories, tags             │
│  Constraints: unique slug per realm                        │
└─────────────────────────────────────────────────────────────┘
```

## Backend Design

### Models

#### WikiPage Model

```go
type WikiPage struct {
    ID              string    `json:"id" gorm:"primaryKey"`
    RealmID         string    `json:"realm_id" gorm:"not null"`
    Title           string    `json:"title" gorm:"not null"`
    Slug            string    `json:"slug" gorm:"not null;uniqueIndex:idx_wiki_realm_slug"`
    Content         string    `json:"content"`
    Summary         string    `json:"summary"`
    Status          string    `json:"status" gorm:"default:draft"`
    Type            string    `json:"type" gorm:"default:article"`
    Protected       bool      `json:"protected" gorm:"default:false"`
    Locked          bool      `json:"locked" gorm:"default:false"`
    Categories      string    `json:"categories"` // JSON array
    Tags            string    `json:"tags"`       // JSON array
    ParentID        *string   `json:"parent_id"`
    MenuOrder       int       `json:"menu_order" gorm:"default:0"`
    RedirectTo      *string   `json:"redirect_to"`
    ViewCount       int       `json:"view_count" gorm:"default:0"`
    EditCount       int       `json:"edit_count" gorm:"default:0"`
    CurrentVersion  int       `json:"current_version" gorm:"default:1"`
    Language        string    `json:"language" gorm:"default:en"`
    TranslationOf   *string   `json:"translation_of"`
    CreatedBy       string    `json:"created_by"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedBy       string    `json:"updated_by"`
    UpdatedAt       time.Time `json:"updated_at"`
    DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
}
```

#### WikiRevision Model

```go
type WikiRevision struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    PageID      string    `json:"page_id" gorm:"not null"`
    Title       string    `json:"title" gorm:"not null"`
    Content     string    `json:"content"`
    Summary     string    `json:"summary"`
    ChangeNote  string    `json:"change_note"`
    Version     int       `json:"version" gorm:"not null"`
    CreatedBy   string    `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
    DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}
```

### Repository Layer

The `WikiRepository` handles all database operations:

- **CRUD Operations**: Create, read, update, delete wiki pages
- **Search & Filter**: Full-text search, category/tag filtering
- **Version Management**: Create and retrieve revisions
- **Statistics**: View counts, edit counts, recent changes

### Service Layer

The `WikiService` provides business logic:

- **Validation**: Page data validation and sanitization
- **Permissions**: Check user permissions for editing/deleting
- **Version Control**: Automatic revision creation on updates
- **Search**: Full-text search with ranking
- **Statistics**: Update view counts, generate analytics

### API Routes

#### Public Routes (No Authentication Required)

```
GET    /api/v1/wiki/pages                    # List published pages
GET    /api/v1/wiki/page/:slug               # Get page by slug
GET    /api/v1/wiki/search                   # Search pages
GET    /api/v1/wiki/category/:category       # Get pages by category
GET    /api/v1/wiki/tag/:tag                 # Get pages by tag
GET    /api/v1/wiki/random                   # Get random page
GET    /api/v1/wiki/recent-changes           # Get recent changes
GET    /api/v1/wiki/page/:slug/history       # Get page history
```

#### Admin Routes (Authentication Required)

```
POST   /api/v1/admin/wiki/pages              # Create new page
PUT    /api/v1/admin/wiki/pages/:id          # Update page
DELETE /api/v1/admin/wiki/pages/:id          # Delete page
POST   /api/v1/admin/wiki/pages/:id/lock     # Lock page
POST   /api/v1/admin/wiki/pages/:id/unlock   # Unlock page
POST   /api/v1/admin/wiki/pages/:id/revisions # Create revision
POST   /api/v1/admin/wiki/pages/:id/revert/:revisionId # Revert to revision
```

## Frontend Design

### Component Architecture

#### Views

1. **WikiView** (`/wiki`)
   - Main wiki dashboard
   - Search functionality
   - Page listing with pagination
   - Category and tag navigation
   - Recent changes sidebar

2. **WikiPageView** (`/wiki/page/:slug`)
   - Display individual wiki pages
   - Markdown rendering
   - Page metadata and actions
   - Related pages sidebar
   - 404 handling with "Create Page" option

3. **WikiEditView** (`/wiki/edit/:slug?`)
   - Create/edit wiki pages
   - Live Markdown preview
   - Form validation
   - Page settings (status, type, protection)

#### Components

1. **WikiSearchBar**
   - Search input with suggestions
   - Real-time search results
   - Search history

2. **WikiSidebar**
   - Categories navigation
   - Tags cloud
   - Recent changes list
   - Special pages (orphaned, wanted, dead-end)

3. **WikiPageList**
   - Grid/list view of pages
   - Pagination controls
   - Action buttons (edit, delete)
   - Empty state handling

4. **WikiPageActions**
   - Edit, delete, lock/unlock buttons
   - Print, share functionality
   - History navigation

5. **WikiNavigation**
   - Breadcrumb navigation
   - Hierarchical page structure

6. **WikiRevisionList**
   - Version history display
   - Compare revisions
   - Revert functionality

### State Management

The `WikiStore` (Pinia) manages:

```typescript
interface WikiState {
  // Data
  pages: WikiPage[]
  currentPage: WikiPage | null
  revisions: WikiRevision[]
  searchResults: WikiPage[]
  
  // UI State
  loading: boolean
  error: string | null
  searchQuery: string
  
  // Pagination
  currentPageNum: number
  pageSize: number
  totalPages: number
  
  // Filters
  filters: {
    status?: string
    type?: string
    category?: string
    tag?: string
  }
}
```

### Routing

```typescript
const routes = [
  {
    path: '/wiki',
    name: 'wiki',
    component: WikiView,
    meta: { requiresAuth: false }
  },
  {
    path: '/wiki/page/:slug',
    name: 'wiki-page',
    component: WikiPageView,
    props: true,
    meta: { requiresAuth: false }
  },
  {
    path: '/wiki/edit/:slug?',
    name: 'wiki-edit',
    component: WikiEditView,
    props: true,
    meta: { requiresAuth: true }
  }
]
```

## API Design

### Request/Response Models

#### CreateWikiPageRequest

```typescript
interface CreateWikiPageRequest {
  realm_id?: string
  title: string
  slug?: string
  content: string
  summary?: string
  status?: 'draft' | 'published' | 'archived' | 'protected'
  type?: 'article' | 'template' | 'category' | 'redirect' | 'stub'
  is_protected?: boolean
  categories?: string[]
  tags?: string[]
  parent_id?: string
  redirect_to?: string
  language?: string
  change_note?: string
}
```

#### WikiPageResponse

```typescript
interface WikiPageResponse {
  id: string
  title: string
  slug: string
  content: string
  summary: string
  status: string
  type: string
  is_protected: boolean
  is_locked: boolean
  categories: string[]
  tags: string[]
  parent_id?: string
  redirect_to?: string
  view_count: number
  edit_count: number
  current_version: number
  language: string
  translation_of?: string
  created_by: string
  created_at: string
  updated_by: string
  updated_at: string
  authenticated?: boolean
  user_id?: string
  can_edit?: boolean
  can_delete?: boolean
}
```

### Error Handling

The API returns consistent error responses:

```json
{
  "error": "Page not found",
  "details": "The requested page does not exist"
}
```

Common HTTP status codes:
- `200` - Success
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (authentication required)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found (page doesn't exist)
- `409` - Conflict (slug already exists)
- `500` - Internal Server Error

## Database Schema

### wiki_pages Table

```sql
CREATE TABLE wiki_pages (
    id TEXT PRIMARY KEY,
    realm_id TEXT NOT NULL,
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    content TEXT,
    summary TEXT,
    status TEXT DEFAULT 'draft',
    type TEXT DEFAULT 'article',
    protected INTEGER DEFAULT 0,
    locked INTEGER DEFAULT 0,
    categories TEXT, -- JSON array
    tags TEXT,       -- JSON array
    parent_id TEXT,
    menu_order INTEGER DEFAULT 0,
    redirect_to TEXT,
    view_count INTEGER DEFAULT 0,
    edit_count INTEGER DEFAULT 0,
    current_version INTEGER DEFAULT 1,
    language TEXT DEFAULT 'en',
    translation_of TEXT,
    created_by TEXT,
    created_at DATETIME,
    updated_by TEXT,
    updated_at DATETIME,
    deleted_at DATETIME
);

-- Indexes
CREATE INDEX idx_wiki_pages_deleted_at ON wiki_pages(deleted_at);
CREATE INDEX idx_wiki_pages_status ON wiki_pages(status);
CREATE INDEX idx_wiki_pages_type ON wiki_pages(type);
CREATE INDEX idx_wiki_pages_parent_id ON wiki_pages(parent_id);
CREATE INDEX idx_wiki_pages_translation_of ON wiki_pages(translation_of);
CREATE UNIQUE INDEX idx_wiki_realm_slug ON wiki_pages(realm_id, slug);
```

### wiki_revisions Table

```sql
CREATE TABLE wiki_revisions (
    id TEXT PRIMARY KEY,
    page_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    summary TEXT,
    change_note TEXT,
    version INTEGER NOT NULL,
    created_by TEXT,
    created_at DATETIME,
    deleted_at DATETIME
);

-- Indexes
CREATE INDEX idx_wiki_revisions_page_id ON wiki_revisions(page_id);
CREATE INDEX idx_wiki_revisions_deleted_at ON wiki_revisions(deleted_at);
```

## Security Model

### Authentication & Authorization

1. **Public Access**: Anyone can view published pages
2. **Authenticated Access**: Required for creating/editing pages
3. **Role-based Permissions**: Admin users can manage all pages
4. **Page Protection**: Individual pages can be protected

### Data Validation

1. **Input Sanitization**: All user input is sanitized
2. **XSS Prevention**: Markdown content is safely rendered
3. **CSRF Protection**: API endpoints use CSRF tokens
4. **SQL Injection Prevention**: Parameterized queries

### Content Security

1. **Markdown Rendering**: Safe HTML generation with `marked` library
2. **File Upload**: Restricted file types and sizes
3. **Rate Limiting**: API rate limiting to prevent abuse

## Markdown Rendering

### Implementation

The frontend uses the `marked` library for Markdown-to-HTML conversion:

```typescript
import { marked } from 'marked'

// Configure marked options
marked.setOptions({
  breaks: true,    // Convert line breaks to <br>
  gfm: true,       // GitHub Flavored Markdown
})

// Render markdown
const html = marked(markdownContent)
```

### Supported Features

- **Headers**: `# H1` through `###### H6`
- **Text Formatting**: `**bold**`, `*italic*`, `***bold italic***`, `~~strikethrough~~`
- **Links**: `[text](url)`, `[text](url "title")`
- **Lists**: Ordered (`1. item`) and unordered (`- item`)
- **Code**: Inline `` `code` `` and code blocks with syntax highlighting
- **Tables**: Full table support with headers
- **Blockquotes**: `> quoted text`
- **Horizontal Rules**: `---` for dividers
- **Images**: `![alt](url)` with responsive styling

### Styling

Custom CSS provides enhanced styling:

```css
.markdown-content a {
  color: #667eea;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: all 0.2s ease;
}

.markdown-content a:hover {
  color: #5a67d8;
  border-bottom-color: #5a67d8;
}

.markdown-content table {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
}

.markdown-content th,
.markdown-content td {
  border: 1px solid #dee2e6;
  padding: 0.75rem;
  text-align: left;
}
```

## Version Control

### Revision System

Every page edit creates a new revision:

1. **Automatic Revisions**: Created on every save
2. **Change Notes**: Optional description of changes
3. **Version Numbering**: Sequential version numbers
4. **Author Tracking**: Track who made each change

### Revision Management

- **View History**: Browse all revisions of a page
- **Compare Revisions**: Side-by-side diff view
- **Revert Changes**: Restore to any previous version
- **Revision Metadata**: Timestamps, authors, change notes

## User Experience

### Page Not Found Handling

When a page doesn't exist (404), the system provides:

1. **Friendly Error Message**: Clear explanation that page doesn't exist
2. **Create Page Option**: Direct link to create the missing page
3. **Authentication Check**: Sign-in prompt for non-authenticated users
4. **Navigation Options**: Back to wiki home, search for similar pages

### Live Preview

The edit interface includes:

1. **Real-time Preview**: Markdown rendered as you type
2. **Side-by-side View**: Editor and preview panels
3. **Syntax Highlighting**: Code blocks with language detection
4. **Form Validation**: Real-time validation feedback

### Search Experience

1. **Instant Search**: Real-time search results
2. **Search Suggestions**: Auto-complete for page titles
3. **Advanced Filters**: By category, tag, status, type
4. **Search History**: Recent searches

## Internationalization

### Supported Languages

- **English** (en) - Default
- **Chinese** (zh) - Simplified Chinese

### Translation Keys

Key translation categories:

```json
{
  "wiki": {
    "title": "Wiki",
    "subtitle": "Collaborative knowledge base",
    "createPage": "Create Page",
    "editPage": "Edit Page",
    "pageNotFound": "Page Not Found",
    "pageNotFoundDesc": "The page \"{slug}\" does not exist yet.",
    "createThisPage": "Create This Page",
    "status": {
      "draft": "Draft",
      "published": "Published",
      "archived": "Archived",
      "protected": "Protected"
    },
    "type": {
      "article": "Article",
      "template": "Template",
      "category": "Category",
      "redirect": "Redirect",
      "stub": "Stub"
    }
  }
}
```

### Locale-specific Features

1. **Date Formatting**: Locale-appropriate date/time display
2. **Number Formatting**: Proper number and currency formatting
3. **Text Direction**: Support for RTL languages
4. **Cultural Considerations**: Appropriate icons and colors

## Future Enhancements

### Planned Features

1. **Advanced Search**: Full-text search with filters and sorting
2. **Page Templates**: Reusable page templates
3. **Collaborative Editing**: Real-time collaborative editing
4. **File Attachments**: Image and document uploads
5. **Page Relationships**: Parent-child page hierarchies
6. **Advanced Permissions**: Granular permission system
7. **API Documentation**: Auto-generated API docs
8. **Export/Import**: Bulk page export and import
9. **Analytics**: Page view analytics and insights
10. **Mobile App**: Native mobile application

### Technical Improvements

1. **Caching**: Redis caching for improved performance
2. **CDN Integration**: Static asset delivery optimization
3. **Search Engine**: Elasticsearch integration
4. **Real-time Updates**: WebSocket support for live updates
5. **Offline Support**: Progressive Web App features
6. **Performance Monitoring**: Application performance monitoring
7. **Automated Testing**: Comprehensive test coverage
8. **Documentation**: API documentation and user guides

## Conclusion

The Wiki module provides a comprehensive, scalable solution for collaborative knowledge management. Its modular architecture, robust security model, and user-friendly interface make it suitable for various use cases, from personal wikis to enterprise knowledge bases.

The design emphasizes:
- **Simplicity**: Easy-to-use interface for content creators
- **Flexibility**: Configurable permissions and page types
- **Scalability**: Efficient database design and caching strategies
- **Security**: Comprehensive input validation and access control
- **Internationalization**: Multi-language support
- **Extensibility**: Plugin architecture for future enhancements
