INITIAL CONTEXT FOR LLM - never change the context-----------------------------
-> THIS SECTION IS A GUIDELINE TO THE LLM CONSIDER BEFORE WORKING IN THIS FILE, DO NOT CHANGE THIS

-> GOES OF THE README FILE:

- This file serves as the technical documentation of the service, providing a comprehensive overview of the codebase
- It should document:
  - Service architecture and design decisions
  - Component structure and relationships
  - API endpoints and interfaces
  - Dependencies and integration points
  - Configuration and deployment details
- This is the primary reference for understanding the technical implementation
- This file should be in sync with the `/TRACKER&MANAGER.md` where development progress and tasks are tracked
- While TRACKER&MANAGER.md focuses on "what" and "when", this file focuses on "how" and "why"

-> CONSIDERER BEFORE UPDATING THIS FILE:

- Never add fictional dates, version numbers, or metrics. Only include real, verified information. If information is not available, mark it as "To be determined" or remove the section.
- The changes in this file need to be incremental or to update informations that you confidentilly have knowlegde, they should not be guesses. If there are questions or uncertanty add comments asking for clarification instead.
- Check the `microservices/docs` folder for comprehensive project details, including architecture, development guidelines, and integration points. This will help in making informed decisions, haver better context and updates to the development plan. Always compare the implementation of this project with the plan described in the docs and whenever there are inconsistancies, add comments.
- Consider structuring this documentation separating the components and describing each one of them and adding sections to how they interact with each other - because this will be very dinamic and updated during the development process it will make clear what to update after each change
- Do not forget to be LLM focus, so because this will be used
- For LLM-specific guidelines and patterns, refer to [LLM Integration Guide](../../docs/llm/README.md)

---

# [Service Name]

-> IMPORTANT: Never add fictional dates, version numbers, or metrics. Only include real, verified information. If information is not available, mark it as "To be determined" or remove the section.

-> CONSIDERER BEFORE UPDATING THIS FILE: The changes in this file need to be incremental, so update tasks as completed our add things that were done but do not remove any of the future plans - only change or insert new things, but do not delete. If something needs to be removed or changed, add a note instead.

-> WHERE TO GET INFORMATION TO IMPROVE THE CONTEXT: Check the `docs` folder for comprehensive project details, including architecture, development guidelines, and integration points. This will help in making informed decisions and updates to the development plan. Remember to update tasks incrementally and document all changes.

-> FOLLOW THE GUIDELINES DEFINED IN THE @TEMPLATES/SERVICES TO UPDATE THE DOCUMENTATION

## Primary Purpose

[Service Name] provides [specific functionality] for the microservices architecture. This service is responsible for [main responsibilities].

## Development Protocol

### 1. Documentation Updates

- Update documentation with all decisions
  - Document design decisions
  - Track implementation details
  - Update cross-references
  - Maintain accuracy
- Maintain cross-references
  - Update related documents
  - Verify link validity
  - Track dependencies
  - Document relationships
- Document implementation details
  - Code structure
  - Configuration
  - Dependencies
  - Integration points
- Track changes
  - Version history
  - Change log
  - Update notes
  - Progress tracking

### 2. Code Development

- Follow development plan
  - Track progress
  - Update tasks
  - Document changes
  - Maintain quality
- Update documentation
  - Code changes
  - Design updates
  - Configuration changes
  - Integration updates
- Maintain tests
  - Unit tests
  - Integration tests
  - Performance tests
  - Security tests
- Track progress
  - Task completion
  - Milestone tracking
  - Issue resolution
  - Quality metrics

### 3. LLM Integration

- Use LLM for code generation
  ```prompt
  # Example: Code Generation
  Generate [component] for [service] that:
  - Implements [requirements]
  - Follows [patterns]
  - Includes [features]
  - Handles [cases]
  ```
- Validate generated code
  - Code review
  - Testing
  - Security check
  - Performance validation
- Document decisions
  - Design choices
  - Implementation details
  - Trade-offs
  - Alternatives
- Update documentation
  - Code changes
  - Design updates
  - Configuration
  - Integration

## Key Documentation References

### Architecture

- [Architecture Overview](../../docs/architecture/README.md)
- [Service Architecture](../../docs/architecture/services/[service-name]-service.md)
- [Security Architecture](../../docs/architecture/overview/security.md)

### Development

- [Development Guide](../../docs/guides/development/guide.md)
- [Testing Guide](../../docs/guides/development/testing/guide.md)
- [Environment Setup](../../docs/guides/development/environment/guide.md)

### API

- [API Specification](../../docs/api/openapi/[service-name]-api.yaml)
- [API Security](../../docs/api/security.md)
- [API Examples](../../docs/api/examples/)

### Operations

- [Monitoring Guide](../../docs/guides/operations/monitoring/guide.md)
- [Logging Guide](../../docs/guides/operations/logging/guide.md)
- [Troubleshooting Guide](../../docs/guides/operations/troubleshooting/guide.md)

## Development Workflow

### 1. Planning

- Review requirements
  - Functional requirements
  - Non-functional requirements
  - Security requirements
  - Performance requirements
- Update documentation
  - Requirements doc
  - Design doc
  - API spec
  - Test plan
- Plan implementation
  - Task breakdown
  - Timeline
  - Dependencies
  - Resources
- Track progress
  - Milestones
  - Tasks
  - Issues
  - Quality

### 2. Implementation

- Follow development plan
  - Code structure
  - Design patterns
  - Best practices
  - Standards
- Update documentation
  - Code changes
  - Design updates
  - Configuration
  - Integration
- Write tests
  - Unit tests
  - Integration tests
  - Performance tests
  - Security tests
- Track changes
  - Version control
  - Change log
  - Update notes
  - Progress

### 3. Review

- Review code
  - Quality check
  - Security review
  - Performance review
  - Best practices
- Update documentation
  - Code changes
  - Design updates
  - Configuration
  - Integration
- Validate tests
  - Test coverage
  - Test quality
  - Test performance
  - Test security
- Track improvements
  - Code quality
  - Documentation
  - Testing
  - Performance

### 4. Deployment

- Prepare deployment
  - Configuration
  - Environment
  - Dependencies
  - Security
- Update documentation
  - Deployment guide
  - Configuration
  - Troubleshooting
  - Monitoring
- Validate deployment
  - Functionality
  - Performance
  - Security
  - Monitoring
- Track status
  - Deployment
  - Monitoring
  - Issues
  - Performance

## Documentation Protocol

### 1. Updates

- Document all decisions
  - Design choices
  - Implementation details
  - Trade-offs
  - Alternatives
- Update cross-references
  - Related docs
  - Dependencies
  - Integration points
  - Dependencies
- Track changes
  - Version history
  - Change log
  - Update notes
  - Progress
- Maintain accuracy
  - Content review
  - Link validation
  - Example updates
  - Code sync

### 2. Reviews

- Review documentation
  - Content accuracy
  - Technical details
  - Examples
  - Cross-references
- Validate accuracy
  - Code sync
  - Configuration
  - Integration
  - Security
- Update references
  - Links
  - Examples
  - Dependencies
  - Integration
- Track improvements
  - Content
  - Structure
  - Examples
  - References

### 3. Maintenance

- Regular updates
  - Code changes
  - Design updates
  - Configuration
  - Integration
- Cross-reference checks
  - Link validation
  - Example updates
  - Dependency tracking
  - Integration points
- Content validation
  - Accuracy
  - Completeness
  - Consistency
  - Quality
- Progress tracking
  - Updates
  - Changes
  - Improvements
  - Quality

## Project Status

- Current Phase: [Phase]
- Status: [Status]
- Next Steps: [Steps]
- Blockers: [Blockers]
- Dependencies: [Dependencies]

## Notes

- Track all decisions
- Update documentation
- Maintain progress
- Document challenges
- Record lessons learned
- Track improvements

## Tasks History

- Initial setup
- Added [feature]
- Updated [section]
- Completed [task]
- Fixed [issue]
- Improved [aspect]
