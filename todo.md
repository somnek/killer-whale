# TODO
> ---------
[ ] search bar:
    - regex/fuzzy search
[ ] sort image by size
[ ] do loading screen cold start
    - do a initial connection to docker, before do anything else
        (for loading screne purpose)
[ ] Learning: show actual command user can use to perform it (filter etc)

# NOTE
> ---------
- main feature should be
    & easy way for user to filter
    & operation on commands


-> 04/04/2023
[ ] status  (by color)
    - restarting (blinking)
    - created, restarting, running, removing, paused, exited, or dead
    * src: [link](https://docs.docker.com/engine/reference/commandline/ps/#filter)

# REF
> --------
- [ticker](https://github.com/achannarasappa/ticker)
- [command tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials/commands/)
- [p.Send](https://github.com/charmbracelet/bubbletea/issues/25#issuecomment-872331380)





# trash 
> ---------------
func listenToDockerEvents(client *docker.Client) tea.Cmd {
	events := make(chan *docker.APIEvents)
	client.AddEventListener(events)
	for event := range events {
		fmt.Println(event)
	}
	return nil
}



