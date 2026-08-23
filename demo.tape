Output demo.gif
Set Width 640
Set Height 480
Set Theme "TokyoNight"
Set FontSize 14
Set FontFamily "Menlo"
Set Padding 2
Set Margin 1
Set Framerate 10

# Launch the TUI
Hide
Type "./opuscriber --interactive --models ./models --model tiny"
Sleep 500ms
Enter
Sleep 2s
Show

# IDLE screen visible
Sleep 2s

# Open file picker by pressing Enter
Enter
Sleep 1s

# Navigate to demo directory
Type "demo"
Sleep 300ms
Enter
Sleep 1s

# Wait for file list to populate
Sleep 500ms

# Tab to select the file, then Enter to confirm
Tab
Sleep 200ms
Enter
Sleep 800ms

# Now in browsing state with file selected, press Tab to start transcribing
Tab
Sleep 500ms

# Wait for transcription to complete
Sleep 6s

# Done screen - show results
Sleep 3s