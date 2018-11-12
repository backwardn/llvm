package asm

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/mewkiz/pkg/diffutil"
	"github.com/mewkiz/pkg/osutil"
)

// words specifies whether to colour words in diff output.
var words bool

func init() {
	flag.BoolVar(&words, "words", false, "colour words in diff output")
	flag.Parse()
}

func TestParseFile(t *testing.T) {
	golden := []struct {
		path string
	}{
		{path: "testdata/inst_binary.ll"},
		{path: "testdata/inst_bitwise.ll"},
		{path: "testdata/inst_vector.ll"},
		{path: "testdata/inst_aggregate.ll"},
		{path: "testdata/inst_memory.ll"},
		{path: "testdata/inst_conversion.ll"},
		{path: "testdata/inst_other.ll"},
		{path: "testdata/terminator.ll"},

		// LLVM Features.
		{path: "testdata/Feature/OperandBundles/adce.ll"},
		{path: "testdata/Feature/OperandBundles/basic-aa-argmemonly.ll"},
		{path: "testdata/Feature/OperandBundles/dse.ll"},
		{path: "testdata/Feature/OperandBundles/early-cse.ll"},
		{path: "testdata/Feature/OperandBundles/function-attrs.ll"},
		{path: "testdata/Feature/OperandBundles/inliner-conservative.ll"},
		{path: "testdata/Feature/OperandBundles/merge-func.ll"},
		{path: "testdata/Feature/OperandBundles/pr26510.ll"},
		{path: "testdata/Feature/OperandBundles/special-state.ll"},
		//{path: "testdata/Feature/alias2.ll"}, // TODO: fix grammar. syntax error at line 12
		//{path: "testdata/Feature/aliases.ll"}, // TODO: fix grammar. syntax error at line 29
		//{path: "testdata/Feature/alignment.ll"}, // TODO: fix grammar. syntax error at line 7
		{path: "testdata/Feature/attributes.ll"},
		{path: "testdata/Feature/basictest.ll"},
		{path: "testdata/Feature/callingconventions.ll"},
		{path: "testdata/Feature/calltest.ll"},
		{path: "testdata/Feature/casttest.ll"},
		{path: "testdata/Feature/cfgstructures.ll"},
		{path: "testdata/Feature/cold.ll"},
		{path: "testdata/Feature/comdat.ll"},
		//{path: "testdata/Feature/constexpr.ll"}, // TODO: re-enable when signed hex integer literals are supported.
		{path: "testdata/Feature/constpointer.ll"},
		{path: "testdata/Feature/const_pv.ll"},
		{path: "testdata/Feature/elf-linker-options.ll"},
		{path: "testdata/Feature/escaped_label.ll"},
		{path: "testdata/Feature/exception.ll"},
		//{path: "testdata/Feature/float.ll"}, // TODO: re-enable when hex half literals are supported.
		{path: "testdata/Feature/fold-fpcast.ll"},
		{path: "testdata/Feature/forwardreftest.ll"},
		{path: "testdata/Feature/fp-intrinsics.ll"},
		{path: "testdata/Feature/global_pv.ll"},
		//{path: "testdata/Feature/globalredefinition3.ll"}, // TODO: figure out how to test. should report error "redefinition of global '@B'"
		{path: "testdata/Feature/global_section.ll"},
		{path: "testdata/Feature/globalvars.ll"},
		{path: "testdata/Feature/indirectcall2.ll"},
		{path: "testdata/Feature/indirectcall.ll"},
		{path: "testdata/Feature/inlineasm.ll"},
		{path: "testdata/Feature/instructions.ll"},
		{path: "testdata/Feature/intrinsic-noduplicate.ll"},
		{path: "testdata/Feature/intrinsics.ll"},
		{path: "testdata/Feature/load_module.ll"},
		{path: "testdata/Feature/md_on_instruction.ll"},
		{path: "testdata/Feature/memorymarkers.ll"},
		{path: "testdata/Feature/metadata.ll"},
		{path: "testdata/Feature/minsize_attr.ll"},
		{path: "testdata/Feature/NamedMDNode2.ll"},
		{path: "testdata/Feature/NamedMDNode.ll"},
		{path: "testdata/Feature/newcasts.ll"},
		{path: "testdata/Feature/optnone.ll"},
		{path: "testdata/Feature/optnone-llc.ll"},
		{path: "testdata/Feature/optnone-opt.ll"},
		{path: "testdata/Feature/packed.ll"},
		{path: "testdata/Feature/packed_struct.ll"},
		{path: "testdata/Feature/paramattrs.ll"},
		{path: "testdata/Feature/ppcld.ll"},
		{path: "testdata/Feature/prefixdata.ll"},
		{path: "testdata/Feature/prologuedata.ll"},
		{path: "testdata/Feature/properties.ll"},
		{path: "testdata/Feature/prototype.ll"},
		{path: "testdata/Feature/recursivetype.ll"},
		{path: "testdata/Feature/seh-nounwind.ll"},
		{path: "testdata/Feature/simplecalltest.ll"},
		{path: "testdata/Feature/smallest.ll"},
		{path: "testdata/Feature/small.ll"},
		{path: "testdata/Feature/sparcld.ll"},
		{path: "testdata/Feature/strip_names.ll"},
		//{path: "testdata/Feature/terminators.ll"}, // TODO: fix grammar. syntax error at line 35
		{path: "testdata/Feature/testalloca.ll"},
		{path: "testdata/Feature/testconstants.ll"},
		{path: "testdata/Feature/testlogical.ll"},
		//{path: "testdata/Feature/testtype.ll"}, // TODO: fix nil pointer dereference
		{path: "testdata/Feature/testvarargs.ll"},
		{path: "testdata/Feature/undefined.ll"},
		{path: "testdata/Feature/unreachable.ll"},
		{path: "testdata/Feature/varargs.ll"},
		{path: "testdata/Feature/varargs_new.ll"},
		{path: "testdata/Feature/vector-cast-constant-exprs.ll"},
		{path: "testdata/Feature/weak_constant.ll"},
		//{path: "testdata/Feature/weirdnames.ll"}, // TODO: re-enable when floats are printed using the same format as Clang.
		//{path: "testdata/Feature/x86ld.ll"}, // TODO: re-enable when floats are printed using the same format as Clang.

		// Coreutils.
		/*
			{path: "testdata/coreutils/[.ll"},
			{path: "testdata/coreutils/b2sum.ll"},
			{path: "testdata/coreutils/base32.ll"},
			{path: "testdata/coreutils/base64.ll"},
			{path: "testdata/coreutils/basename.ll"},
			{path: "testdata/coreutils/cat.ll"},
			{path: "testdata/coreutils/chcon.ll"},
			{path: "testdata/coreutils/chgrp.ll"},
			{path: "testdata/coreutils/chmod.ll"},
			{path: "testdata/coreutils/chown.ll"},
			{path: "testdata/coreutils/chroot.ll"},
			{path: "testdata/coreutils/cksum.ll"},
			{path: "testdata/coreutils/comm.ll"},
			{path: "testdata/coreutils/cp.ll"},
			{path: "testdata/coreutils/csplit.ll"},
			{path: "testdata/coreutils/cut.ll"},
			{path: "testdata/coreutils/date.ll"},
			{path: "testdata/coreutils/dd.ll"},
			{path: "testdata/coreutils/df.ll"},
			{path: "testdata/coreutils/dir.ll"},
			{path: "testdata/coreutils/dircolors.ll"},
			{path: "testdata/coreutils/dirname.ll"},
			{path: "testdata/coreutils/du.ll"},
			{path: "testdata/coreutils/echo.ll"},
			{path: "testdata/coreutils/env.ll"},
			{path: "testdata/coreutils/expand.ll"},
			{path: "testdata/coreutils/expr.ll"},
			{path: "testdata/coreutils/factor.ll"},
			{path: "testdata/coreutils/false.ll"},
			{path: "testdata/coreutils/fmt.ll"},
			{path: "testdata/coreutils/fold.ll"},
			{path: "testdata/coreutils/getlimits.ll"},
			{path: "testdata/coreutils/ginstall.ll"},
			{path: "testdata/coreutils/groups.ll"},
			{path: "testdata/coreutils/head.ll"},
			{path: "testdata/coreutils/hostid.ll"},
			{path: "testdata/coreutils/id.ll"},
			{path: "testdata/coreutils/join.ll"},
			{path: "testdata/coreutils/kill.ll"},
			{path: "testdata/coreutils/link.ll"},
			{path: "testdata/coreutils/ln.ll"},
			{path: "testdata/coreutils/logname.ll"},
			{path: "testdata/coreutils/ls.ll"},
			{path: "testdata/coreutils/make-prime-list.ll"},
			{path: "testdata/coreutils/md5sum.ll"},
			{path: "testdata/coreutils/mkdir.ll"},
			{path: "testdata/coreutils/mkfifo.ll"},
			{path: "testdata/coreutils/mknod.ll"},
			{path: "testdata/coreutils/mktemp.ll"},
			{path: "testdata/coreutils/mv.ll"},
			{path: "testdata/coreutils/nice.ll"},
			{path: "testdata/coreutils/nl.ll"},
			{path: "testdata/coreutils/nohup.ll"},
			{path: "testdata/coreutils/nproc.ll"},
			{path: "testdata/coreutils/numfmt.ll"},
			{path: "testdata/coreutils/od.ll"},
			{path: "testdata/coreutils/paste.ll"},
			{path: "testdata/coreutils/pathchk.ll"},
			{path: "testdata/coreutils/pinky.ll"},
			{path: "testdata/coreutils/pr.ll"},
			{path: "testdata/coreutils/printenv.ll"},
			{path: "testdata/coreutils/printf.ll"},
			{path: "testdata/coreutils/ptx.ll"},
			{path: "testdata/coreutils/pwd.ll"},
			{path: "testdata/coreutils/readlink.ll"},
			{path: "testdata/coreutils/realpath.ll"},
			{path: "testdata/coreutils/rm.ll"},
			{path: "testdata/coreutils/rmdir.ll"},
			{path: "testdata/coreutils/runcon.ll"},
			{path: "testdata/coreutils/seq.ll"},
			{path: "testdata/coreutils/sha1sum.ll"},
			{path: "testdata/coreutils/sha224sum.ll"},
			{path: "testdata/coreutils/sha256sum.ll"},
			{path: "testdata/coreutils/sha384sum.ll"},
			{path: "testdata/coreutils/sha512sum.ll"},
			{path: "testdata/coreutils/shred.ll"},
			{path: "testdata/coreutils/shuf.ll"},
			{path: "testdata/coreutils/sleep.ll"},
			{path: "testdata/coreutils/sort.ll"},
			{path: "testdata/coreutils/split.ll"},
			{path: "testdata/coreutils/stat.ll"},
			{path: "testdata/coreutils/stdbuf.ll"},
			{path: "testdata/coreutils/stty.ll"},
			{path: "testdata/coreutils/sum.ll"},
			{path: "testdata/coreutils/sync.ll"},
			{path: "testdata/coreutils/tac.ll"},
			{path: "testdata/coreutils/tail.ll"},
			{path: "testdata/coreutils/tee.ll"},
			{path: "testdata/coreutils/test.ll"},
			{path: "testdata/coreutils/timeout.ll"},
			{path: "testdata/coreutils/touch.ll"},
			{path: "testdata/coreutils/tr.ll"},
			{path: "testdata/coreutils/true.ll"},
			{path: "testdata/coreutils/truncate.ll"},
			{path: "testdata/coreutils/tsort.ll"},
			{path: "testdata/coreutils/tty.ll"},
			{path: "testdata/coreutils/uname.ll"},
			{path: "testdata/coreutils/unexpand.ll"},
			{path: "testdata/coreutils/uniq.ll"},
			{path: "testdata/coreutils/unlink.ll"},
			{path: "testdata/coreutils/uptime.ll"},
			{path: "testdata/coreutils/users.ll"},
			{path: "testdata/coreutils/vdir.ll"},
			{path: "testdata/coreutils/wc.ll"},
			{path: "testdata/coreutils/who.ll"},
			{path: "testdata/coreutils/whoami.ll"},
			{path: "testdata/coreutils/yes.ll"},
		*/

		// SQLite.
		//{path: "testdata/sqlite/shell.ll"},
	}
	for _, g := range golden {
		m, err := ParseFile(g.path)
		if err != nil {
			t.Errorf("unable to parse %q into AST; %+v", g.path, err)
			continue
		}
		path := g.path
		if osutil.Exists(g.path + ".golden") {
			path = g.path + ".golden"
		}
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf("unable to read %q; %+v", path, err)
			continue
		}
		want := string(buf)
		got := m.String()
		if want != got {
			if err := diffutil.Diff(want, got, words, filepath.Base(path)); err != nil {
				panic(err)
			}
			t.Errorf("module mismatch %q; expected `%s`, got `%s`", path, want, got)
			continue
		}
	}
}