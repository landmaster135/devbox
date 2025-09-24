Implement some functionality as instructed by the user.

<detailed_sequence_of_steps>

# Implement functionality

1. Retrieve the contents of `.agents/instructions.md` located in the root of the current project.
2. If the retrieved instructions are empty, ask the user to write instructions in `.agents/instructions.md` and obtain their feedback.
3. If the retrieved instructions refer to something that was carried out recently, ask the user again whether to proceed with the implementation and obtain their feedback.
4. After receiving the user’s feedback, determine whether there are any further questions you need to ask.
5. Once no further questions are needed, proceed with the implementation.
