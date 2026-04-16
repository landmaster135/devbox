---
name: implementing-slowly
description: Work any task as instructed by the user.
---

# Detailed Sequence of Steps

1. Retrieve the latest contents of the instruction document in the current project.
2. If the retrieved instructions are empty, ask the user to write instructions in the instruction document and obtain their feedback.
3. If the retrieved instructions refer to something that was carried out recently, ask the user again whether to proceed with the task and obtain their feedback.
4. After receiving the user’s feedback, determine whether there are any further questions you need to ask.
5. Once no further questions are needed, proceed with the task.
