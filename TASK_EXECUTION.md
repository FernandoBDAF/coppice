# Task Execution and Verification Guide

## Task Execution Flow

### 1. Initial Task Setup

- Generate prompt using `LLM.md` template
- Review and validate prompt completeness
- Start new chat with Cursor using generated prompt

### 2. Iterative Execution Process

#### Phase 1: Planning

1. Review task requirements
2. Break down into sub-tasks
3. Create implementation checklist
4. Identify potential challenges

#### Phase 2: Implementation

1. Execute each sub-task
2. Document changes made
3. Track progress against checklist
4. Handle any blockers

#### Phase 3: Verification

1. Review changes against requirements
2. Run tests and validations
3. Update documentation
4. Final review

## Task Completion Checklist

### Implementation Verification

- [ ] All requirements implemented
- [ ] Code follows standards
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No regression issues

### Documentation Updates

- [ ] README.md updated
- [ ] CONTEXT.MD updated
- [ ] INTERFACE.MD updated
- [ ] MANAGER.md updated
- [ ] TRACKER.MD updated

### Quality Checks

- [ ] Code review completed
- [ ] Tests passing
- [ ] Linting passed
- [ ] Build successful
- [ ] No security issues

## Task Iteration Protocol

### 1. Progress Tracking

```markdown
## Current Status

- Completed: [list of completed items]
- In Progress: [current task]
- Pending: [remaining tasks]

## Next Steps

1. [immediate next action]
2. [following action]
3. [subsequent action]

## Blockers/Questions

- [list any blockers or questions]
```

### 2. Change Documentation

```markdown
## Changes Made

- File: [filename]
  - Changes: [description]
  - Impact: [effect]
  - Verification: [how verified]

## Documentation Updates

- File: [filename]
  - Section: [section name]
  - Updates: [description]
```

### 3. Verification Steps

```markdown
## Verification Checklist

- [ ] Implementation matches requirements
- [ ] All tests passing
- [ ] Documentation updated
- [ ] No regression issues
- [ ] Performance acceptable
```

## Task Completion Protocol

### 1. Final Review

- Verify all requirements met
- Check all documentation updated
- Ensure all tests passing
- Validate no regression issues

### 2. Documentation Updates

- Update TRACKER.MD with completion status
- Update relevant documentation files
- Add any new patterns or learnings

### 3. Handover

- Document any special considerations
- Note any future improvements
- Update task status in TRACKER.MD

## Example Task Iteration

```markdown
# Task Progress Update

## Current Status

- Completed:
  - Initial setup
  - Basic implementation
- In Progress:
  - Advanced features
- Pending:
  - Testing
  - Documentation

## Next Steps

1. Complete advanced features
2. Write tests
3. Update documentation

## Changes Made

- File: service/example.js
  - Changes: Added new feature
  - Impact: Improved functionality
  - Verification: Unit tests added

## Verification Checklist

- [x] Basic implementation
- [ ] Advanced features
- [ ] Tests
- [ ] Documentation
```

## Task Completion Template

```markdown
# Task Completion Report

## Implementation Summary

- [Brief summary of what was implemented]

## Changes Made

- [List of all changes]

## Documentation Updates

- [List of documentation updates]

## Verification Results

- [Test results]
- [Performance metrics]
- [Security checks]

## Future Considerations

- [Any future improvements]
- [Potential optimizations]
- [Maintenance notes]
```

## Using This Guide

1. **Start New Task**

   - Generate prompt using LLM.md
   - Create new chat with Cursor
   - Begin implementation

2. **Track Progress**

   - Use Task Iteration Protocol
   - Update status regularly
   - Document changes

3. **Verify Completion**

   - Use Task Completion Protocol
   - Complete all checklists
   - Update documentation

4. **Handover**
   - Complete Task Completion Report
   - Update TRACKER.MD
   - Document any special considerations
