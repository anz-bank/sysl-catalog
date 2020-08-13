package catalog

import "fmt"

/* String returns a string of all of the non pointer fields; mainly to be used with require.Equal*/
func (p *Generator) String() string {
	return fmt.Sprint(
		p.FilesToCreate,
		p.MermaidFilesToCreate,
		p.RedocFilesToCreate,
		p.GeneratedFiles,
		p.SourceFileName,
		p.ProjectTitle,
		p.ImageDest,
		p.Format,
		p.Ext,
		p.OutputFileName,
		p.PlantumlService,
		p.StartTemplateIndex,
		p.FilterPackage,
		p.CustomTemplate,
		p.LiveReload,
		p.ImageTags,
		p.DisableCss,
		p.DisableImages,
		p.Mermaid,
		p.Redoc,
		p.Fs,
		p.errs,
		p.CurrentDir,
		p.TempDir,
		p.Title,
		p.OutputDir,
		p.Links,
		p.Server,
	)
}
