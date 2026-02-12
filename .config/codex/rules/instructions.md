# What is important commonly
- Reply in Japanese.
- When the context is unclear, confirm with the user.
- At first, before beginning work, report on the current situation and the next tasks to be performed.
- When you are tasked with work, search for any MCP tools that you could potentially utilize.
- Before using the MCP server, clearly state the server name and tool name before requesting permission to use it.
- When asked a question by a user, respond with the cause of the problem and the solution before editing the file.
- Before each task, report the situation and the content of the work you're undertaking.

# What is important for coding

## Common things
- When encountering modules or packages with unclear specifications, use the MCP tool to search and investigate their specifications.
- When adding new functions or methods, write them after the final line of similar function groups within the file.
- Create functions and methods in small units whenever possible.
- Don't complete the analysis prematurely, continue analyzing even if you think you found a solution.
- Keep each file's line count under 700 lines. Before adding test code, count the number of lines in the file you're planning to add to using the shell command: "awk 'END { print NR }' ". Then, determine whether that number exceeds 700 using the MCP tool. If it's less than 700, append to the existing test file. If it's more than 700, create a new file in the same directory and append to that file.
- Avoid using hard coding as much as possible.
- Don’t forget to add logging.

## Test Driven Development (TDD)
- When implementing features requested by users, always implement test code for those features as well.
- After implementing features requested by users, confirm with the user before implementing the test code for those features.
- When implementing test code, keep the Red-Green-Refactor cycle in mind.
- Test methods are implemented based on the Arrange-Act-Assert pattern.
- Once the implemented test code passes normally, check the coverage and report it to the user. Then, improve the coverage.

## SOLID principles
Be mindful of SOLID principles. SOLID acronym stands for five design principles that help make software more maintainable and scalable:
- Single Responsibility Principle
- Open/Closed Principle
- Liskov Substitution Principle
- Interface Segregation Principle
- Dependency Inversion Principle
